package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"

	"github.com/dexguitar/spacecraftory/order/internal/config"
	orderMetrics "github.com/dexguitar/spacecraftory/order/internal/metrics"
	customMiddleware "github.com/dexguitar/spacecraftory/order/internal/middleware"
	"github.com/dexguitar/spacecraftory/platform/pkg/closer"
	"github.com/dexguitar/spacecraftory/platform/pkg/logger"
	"github.com/dexguitar/spacecraftory/platform/pkg/metrics"
	httpAuth "github.com/dexguitar/spacecraftory/platform/pkg/middleware/http"
	"github.com/dexguitar/spacecraftory/platform/pkg/migrator"
	pgMigrator "github.com/dexguitar/spacecraftory/platform/pkg/migrator/pg"
	"github.com/dexguitar/spacecraftory/platform/pkg/tracing"
	orderV1 "github.com/dexguitar/spacecraftory/shared/pkg/openapi/order/v1"
)

type App struct {
	diContainer *diContainer
	httpServer  *http.Server
	migrator    migrator.Migrator
}

func New(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	// Канал для ошибок от компонентов
	errCh := make(chan error, 2)

	// Контекст для остановки всех горутин
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Консьюмер
	go func() {
		if err := a.runConsumer(ctx); err != nil {
			errCh <- fmt.Errorf("consumer crashed: %w", err)
		}
	}()

	// HTTP сервер
	go func() {
		if err := a.runHTTPServer(ctx); err != nil {
			errCh <- fmt.Errorf("http server crashed: %w", err)
		}
	}()

	// Ожидание либо ошибки, либо завершения контекста (например, сигнал SIGINT/SIGTERM)
	select {
	case <-ctx.Done():
		logger.Info(ctx, "Shutdown signal received")
	case err := <-errCh:
		logger.Error(ctx, "Component crashed, shutting down", zap.Error(err))
		// Триггерим cancel, чтобы остановить второй компонент
		cancel()
		// Дождись завершения всех задач (если есть graceful shutdown внутри)
		<-ctx.Done()
		return err
	}

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDI,
		a.initLogger,
		a.initMetrics,
		a.initTracing,
		a.initCloser,
		a.initMigrator,
		a.initHTTPServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initDI(_ context.Context) error {
	a.diContainer = NewDiContainer()
	return nil
}

func (a *App) initLogger(_ context.Context) error {
	cfg := config.AppConfig().Logger
	if cfg.OtelEndpoint() != "" {
		return logger.InitWithOTLP(
			cfg.Level(),
			cfg.AsJson(),
			cfg.OtelEndpoint(),
			cfg.ServiceName(),
			"dev",
		)
	}

	return logger.Init(cfg.Level(), cfg.AsJson())
}

func (a *App) initMetrics(ctx context.Context) error {
	cfg := config.AppConfig().Metrics
	if cfg.CollectorEndpoint() == "" {
		return nil // Metrics disabled
	}

	if err := metrics.InitProvider(ctx, cfg); err != nil {
		return fmt.Errorf("failed to init metrics provider: %w", err)
	}

	if err := orderMetrics.InitMetrics(); err != nil {
		return fmt.Errorf("failed to init order metrics: %w", err)
	}

	logger.Info(ctx, "📊 Metrics initialized")
	return nil
}

func (a *App) initTracing(ctx context.Context) error {
	cfg := config.AppConfig().Tracing
	if cfg.CollectorEndpoint() == "" {
		return nil // Tracing disabled
	}

	if err := tracing.InitTracer(ctx, cfg); err != nil {
		return fmt.Errorf("failed to init tracer: %w", err)
	}

	logger.Info(ctx, "🔍 Tracing initialized")
	return nil
}

func (a *App) initCloser(_ context.Context) error {
	closer.SetLogger(logger.Logger())
	closer.AddNamed("Tracing", tracing.ShutdownTracer)
	closer.AddNamed("Metrics", metrics.Shutdown)
	closer.AddNamed("Logger", func(ctx context.Context) error {
		return logger.Close()
	})
	return nil
}

func (a *App) initMigrator(ctx context.Context) error {
	dbURI := config.AppConfig().Postgres.Address()

	conn, err := pgx.Connect(ctx, dbURI)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	err = conn.Ping(ctx)
	if err != nil {
		if closeErr := conn.Close(ctx); closeErr != nil {
			logger.Error(ctx, "❌ failed to close database connection", zap.Error(closeErr))
		}
		return fmt.Errorf("database is unavailable: %w", err)
	}

	migrationsDir := config.AppConfig().Postgres.MigrationDirectory()
	sqlDB := stdlib.OpenDB(*conn.Config().Copy())
	a.migrator = pgMigrator.NewMigrator(sqlDB, migrationsDir)

	logger.Info(ctx, "🔄 Running database migrations...")
	err = a.migrator.Up(ctx)
	if err != nil {
		if closeErr := conn.Close(ctx); closeErr != nil {
			logger.Error(ctx, "❌ failed to close database connection", zap.Error(closeErr))
		}
		if closeErr := sqlDB.Close(); closeErr != nil {
			logger.Error(ctx, "❌ failed to close database connection", zap.Error(closeErr))
		}
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	logger.Info(ctx, "✅ Database migrations completed")

	if closeErr := sqlDB.Close(); closeErr != nil {
		logger.Error(ctx, "❌ failed to close database connection", zap.Error(closeErr))
	}
	if closeErr := conn.Close(ctx); closeErr != nil {
		logger.Error(ctx, "❌ failed to close database connection", zap.Error(closeErr))
	}

	return nil
}

func (a *App) initHTTPServer(ctx context.Context) error {
	orderServer, err := orderV1.NewServer(a.diContainer.OrderV1API(ctx))
	if err != nil {
		return fmt.Errorf("failed to create OpenAPI server: %w", err)
	}

	mux := chi.NewRouter()

	authMiddleware := httpAuth.NewAuthMiddleware(a.diContainer.IAMGRPCClient(ctx))
	mux.Use(authMiddleware.Handle)
	mux.Use(customMiddleware.RequestLogger)
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.Timeout(10 * time.Second))

	mux.Mount("/", orderServer)

	a.httpServer = &http.Server{
		Addr:              config.AppConfig().HTTP.Address(),
		Handler:           mux,
		ReadHeaderTimeout: config.AppConfig().HTTP.ReadTimeout(),
	}

	closer.AddNamed("HTTP server", func(ctx context.Context) error {
		shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		if err := a.httpServer.Shutdown(shutdownCtx); err != nil {
			if !errors.Is(err, http.ErrServerClosed) && !errors.Is(err, net.ErrClosed) {
				return err
			}
		}

		return nil
	})

	return nil
}

func (a *App) runHTTPServer(ctx context.Context) error {
	addr := config.AppConfig().HTTP.Address()
	log.Println("addr", addr)
	logger.Info(ctx, fmt.Sprintf("🚀 Order HTTP server listening on %s", addr))
	logger.Info(ctx, fmt.Sprintf("📚 API available at: http://%s/api/v1/orders", addr))

	err := a.httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (a *App) runConsumer(ctx context.Context) error {
	err := a.diContainer.OrderConsumerService(ctx).RunConsumer(ctx)
	if err != nil {
		return err
	}

	return nil
}

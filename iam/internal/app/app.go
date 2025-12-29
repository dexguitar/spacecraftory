package app

import (
	"context"
	"errors"
	"fmt"
	"net"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/dexguitar/spacecraftory/iam/internal/config"
	"github.com/dexguitar/spacecraftory/platform/pkg/closer"
	"github.com/dexguitar/spacecraftory/platform/pkg/grpc/health"
	"github.com/dexguitar/spacecraftory/platform/pkg/logger"
	"github.com/dexguitar/spacecraftory/platform/pkg/migrator"
	pgMigrator "github.com/dexguitar/spacecraftory/platform/pkg/migrator/pg"
	authV1 "github.com/dexguitar/spacecraftory/shared/pkg/proto/auth/v1"
	userV1 "github.com/dexguitar/spacecraftory/shared/pkg/proto/user/v1"
)

type App struct {
	diContainer *diContainer
	grpcServer  *grpc.Server
	migrator    migrator.Migrator
	listener    net.Listener
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
	return a.runGRPCServer(ctx)
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDI,
		a.initLogger,
		a.initCloser,
		a.initListener,
		a.initGRPCServer,
		a.initMigrator,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

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

func (a *App) initCloser(_ context.Context) error {
	closer.SetLogger(logger.Logger())
	closer.AddNamed("Logger", func(ctx context.Context) error {
		return logger.Close()
	})
	return nil
}

func (a *App) initListener(_ context.Context) error {
	listener, err := net.Listen("tcp", config.AppConfig().IAMGRPC.Address())
	if err != nil {
		return err
	}
	closer.AddNamed("TCP listener", func(ctx context.Context) error {
		lerr := listener.Close()
		if lerr != nil && !errors.Is(lerr, net.ErrClosed) {
			return lerr
		}
		return nil
	})

	a.listener = listener

	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer(grpc.Creds(insecure.NewCredentials()))
	closer.AddNamed("gRPC server", func(ctx context.Context) error {
		a.grpcServer.GracefulStop()
		return nil
	})

	reflection.Register(a.grpcServer)

	// Register health service for health checks
	health.RegisterService(a.grpcServer)

	// Register IAM services
	userV1.RegisterUserServiceServer(a.grpcServer, a.diContainer.UserV1API(ctx))
	authV1.RegisterAuthServiceServer(a.grpcServer, a.diContainer.AuthV1API(ctx))

	// Register Envoy External Authorization service
	authv3.RegisterAuthorizationServer(a.grpcServer, a.diContainer.ExtAuthzAPI(ctx))

	return nil
}

func (a *App) runGRPCServer(ctx context.Context) error {
	logger.Info(ctx, fmt.Sprintf("🚀 IAM gRPC server listening on %s", config.AppConfig().IAMGRPC.Address()))

	err := a.grpcServer.Serve(a.listener)
	if err != nil {
		return err
	}

	return nil
}

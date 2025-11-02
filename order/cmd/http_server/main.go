package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderApi "github.com/dexguitar/spacecraftory/order/internal/api/order/v1"
	inventoryClient "github.com/dexguitar/spacecraftory/order/internal/client/grpc/inventory/v1"
	paymentClient "github.com/dexguitar/spacecraftory/order/internal/client/grpc/payment/v1"
	"github.com/dexguitar/spacecraftory/order/internal/config"
	customMiddleware "github.com/dexguitar/spacecraftory/order/internal/middleware"
	"github.com/dexguitar/spacecraftory/order/internal/migrator"
	orderRepository "github.com/dexguitar/spacecraftory/order/internal/repository/order"
	orderService "github.com/dexguitar/spacecraftory/order/internal/service/order"
	orderV1 "github.com/dexguitar/spacecraftory/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/dexguitar/spacecraftory/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/dexguitar/spacecraftory/shared/pkg/proto/payment/v1"
)

const configPath = "./deploy/compose/order/.env"

func main() {
	err := config.Load(configPath)
	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}

	inventoryConn, err := grpc.NewClient(
		fmt.Sprintf(":%s", config.AppConfig().GRPCClient.InventoryAddress()),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("Failed to connect to inventory service: %v", err)
		return
	}
	defer func() {
		if cerr := inventoryConn.Close(); cerr != nil {
			log.Printf("Failed to close inventory connection: %v", cerr)
		}
	}()

	paymentConn, err := grpc.NewClient(
		fmt.Sprintf(":%s", config.AppConfig().GRPCClient.PaymentAddress()),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("Failed to connect to payment service: %v", err)
		return
	}
	defer func() {
		if cerr := paymentConn.Close(); cerr != nil {
			log.Printf("Failed to close payment connection: %v", cerr)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbURI := config.AppConfig().Postgres.Address()

	conn, err := pgx.Connect(ctx, dbURI)
	if err != nil {
		log.Printf("failed to connect to database: %v\n", err)
		return
	}
	defer func() {
		if cerr := conn.Close(ctx); cerr != nil {
			log.Printf("failed to close connection: %v", cerr)
		}
	}()

	err = conn.Ping(ctx)
	if err != nil {
		log.Printf("Database is unavailable: %v\n", err)
		return
	}

	migrationsDir := config.AppConfig().Postgres.MigrationDirectory()
	migratorRunner := migrator.NewMigrator(stdlib.OpenDB(*conn.Config().Copy()), migrationsDir)

	err = migratorRunner.Up()
	if err != nil {
		log.Printf("failed to run migrations: %v\n", err)
		return
	}

	pool, err := pgxpool.New(ctx, dbURI)
	if err != nil {
		log.Printf("failed to create pool: %v\n", err)
		return
	}
	defer pool.Close()

	repo := orderRepository.NewOrderRepository(pool)

	inventoryGrpcClient := inventoryV1.NewInventoryServiceClient(inventoryConn)
	invClient := inventoryClient.NewInventoryClient(inventoryGrpcClient)

	paymentGrpcClient := paymentV1.NewPaymentServiceClient(paymentConn)
	paymentClient := paymentClient.NewPaymentClient(paymentGrpcClient)

	service := orderService.NewService(repo, invClient, paymentClient)
	api := orderApi.NewAPI(service)

	orderServer, err := orderV1.NewServer(api)
	if err != nil {
		log.Printf("Failed to create OpenAPI server: %v", err)
		return
	}

	mux := chi.NewRouter()

	mux.Use(customMiddleware.RequestLogger)
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.Timeout(10 * time.Second))

	mux.Mount("/", orderServer)

	httpServer := &http.Server{
		Addr:              fmt.Sprintf(":%s", config.AppConfig().HTTP.Address()),
		Handler:           mux,
		ReadHeaderTimeout: config.AppConfig().HTTP.ReadTimeout(),
	}

	go func() {
		addr := config.AppConfig().HTTP.Address()
		log.Printf("🚀 Order HTTP server listening on %s\n", addr)
		log.Printf("📚 API available at: %s/api/v1/orders\n", addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to serve Order HTTP: %v\n", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down Order server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), config.AppConfig().HTTP.ReadTimeout())
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("Order HTTP server shutdown error: %v", err)
	}
	log.Println("✅ Order HTTP server stopped")
}

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
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderApi "github.com/dexguitar/spacecraftory/order/internal/api/order/v1"
	inventoryClient "github.com/dexguitar/spacecraftory/order/internal/client/grpc/inventory/v1"
	paymentClient "github.com/dexguitar/spacecraftory/order/internal/client/grpc/payment/v1"
	customMiddleware "github.com/dexguitar/spacecraftory/order/internal/middleware"
	orderRepository "github.com/dexguitar/spacecraftory/order/internal/repository/order"
	orderService "github.com/dexguitar/spacecraftory/order/internal/service/order"
	orderV1 "github.com/dexguitar/spacecraftory/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/dexguitar/spacecraftory/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/dexguitar/spacecraftory/shared/pkg/proto/payment/v1"
)

const (
	httpPort             = "8080"
	inventoryServiceAddr = "50051"
	paymentServiceAddr   = "50052"

	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

func main() {
	inventoryConn, err := grpc.NewClient(
		fmt.Sprintf(":%s", inventoryServiceAddr),
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
		fmt.Sprintf(":%s", paymentServiceAddr),
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

	repo := orderRepository.NewOrderRepository()

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
		Addr:              fmt.Sprintf(":%s", httpPort),
		Handler:           mux,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	go func() {
		log.Printf("🚀 Order HTTP server listening on port %s\n", httpPort)
		log.Printf("📚 API available at: http://localhost:%s/api/v1/orders\n", httpPort)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to serve Order HTTP: %v\n", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down Order server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("Order HTTP server shutdown error: %v", err)
	}
	log.Println("✅ Order HTTP server stopped")
}

package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	inventoryApi "github.com/dexguitar/spacecraftory/inventory/internal/api/inventory/v1"
	"github.com/dexguitar/spacecraftory/inventory/internal/interceptor"
	inventoryRepository "github.com/dexguitar/spacecraftory/inventory/internal/repository/inventory"
	inventoryService "github.com/dexguitar/spacecraftory/inventory/internal/service/inventory"
	inventoryV1 "github.com/dexguitar/spacecraftory/shared/pkg/proto/inventory/v1"
)

const (
	grpcPort = 50051
)

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}
	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("failed to close listener: %v\n", cerr)
		}
	}()

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(interceptor.ValidationInterceptor()),
	)

	// Enable reflection for debugging
	reflection.Register(s)

	ctx := context.Background()

	err = godotenv.Load(".env")
	if err != nil {
		log.Printf("failed to load .env file: %v\n", err)
		return
	}

	dbURI := os.Getenv("MONGO_URI")

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbURI))
	if err != nil {
		log.Printf("failed to connect to database: %v\n", err)
		return
	}
	defer func() {
		if cerr := client.Disconnect(ctx); cerr != nil {
			log.Printf("failed to disconnect from database: %v\n", cerr)
		}
	}()

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Printf("failed to ping database: %v\n", err)
		return
	}

	db := client.Database(os.Getenv("MONGO_INITDB_DATABASE"))

	repo := inventoryRepository.NewInventoryRepository(db)
	service := inventoryService.NewService(repo)
	api := inventoryApi.NewAPI(service)
	inventoryV1.RegisterInventoryServiceServer(s, api)
	// Start gRPC server
	go func() {
		log.Printf("🚀 Inventory gRPC server listening on port %d\n", grpcPort)
		err = s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve Inventory gRPC: %v\n", err)
			return
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down Inventory servers...")
	s.GracefulStop()
	log.Println("✅ Inventory gRPC server stopped")
}

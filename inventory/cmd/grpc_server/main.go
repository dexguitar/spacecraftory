package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	inventoryApi "github.com/dexguitar/spacecraftory/inventory/internal/api/inventory/v1"
	"github.com/dexguitar/spacecraftory/inventory/internal/config"
	"github.com/dexguitar/spacecraftory/inventory/internal/interceptor"
	inventoryRepository "github.com/dexguitar/spacecraftory/inventory/internal/repository/inventory"
	inventoryService "github.com/dexguitar/spacecraftory/inventory/internal/service/inventory"
	inventoryV1 "github.com/dexguitar/spacecraftory/shared/pkg/proto/inventory/v1"
)

const configPath = "./deploy/compose/inventory/.env"

func main() {
	err := config.Load(configPath)
	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}

	// Create MongoDB client
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(config.AppConfig().Mongo.URI()))
	if err != nil {
		log.Printf("failed to connect to MongoDB: %v\n", err)
		return
	}
	defer func() {
		if cerr := mongoClient.Disconnect(context.Background()); cerr != nil {
			log.Printf("failed to disconnect from MongoDB: %v\n", cerr)
		}
	}()

	// Verify MongoDB connection
	err = mongoClient.Ping(ctx, nil)
	if err != nil {
		log.Printf("failed to ping MongoDB: %v\n", err)
		return
	}
	log.Println("✅ Connected to MongoDB")

	lis, err := net.Listen("tcp", config.AppConfig().InventoryGRPC.Address())
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

	db := mongoClient.Database(config.AppConfig().Mongo.DatabaseName())
	repo := inventoryRepository.NewInventoryRepository(db)
	service := inventoryService.NewService(repo)
	api := inventoryApi.NewAPI(service)
	inventoryV1.RegisterInventoryServiceServer(s, api)

	// Start gRPC server
	go func() {
		log.Printf("🚀 Inventory gRPC server listening on %s\n", config.AppConfig().InventoryGRPC.Address())
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

package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	inventoryApi "github.com/dexguitar/spacecraftory/inventory/internal/api/inventory/v1"
	"github.com/dexguitar/spacecraftory/inventory/internal/interceptor"
	inventoryRepository "github.com/dexguitar/spacecraftory/inventory/internal/repository/inventory"
	"github.com/dexguitar/spacecraftory/inventory/internal/repository/utils"
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

	repo := inventoryRepository.NewInventoryRepository()

	repo.InitializeMockData(utils.GenerateMockParts())

	service := inventoryService.NewService(repo)
	api := inventoryApi.NewAPI(service)

	inventoryV1.RegisterInventoryServiceServer(s, api)

	// Enable reflection for debugging
	reflection.Register(s)

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

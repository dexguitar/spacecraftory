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

	paymentApi "github.com/dexguitar/spacecraftory/payment/internal/api/payment/v1"
	"github.com/dexguitar/spacecraftory/payment/internal/interceptor"
	paymentRepository "github.com/dexguitar/spacecraftory/payment/internal/repository/payment"
	paymentService "github.com/dexguitar/spacecraftory/payment/internal/service/payment"
	paymentV1 "github.com/dexguitar/spacecraftory/shared/pkg/proto/payment/v1"
)

const (
	grpcPort = 50052
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

	// Create gRPC server
	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(interceptor.ValidationInterceptor()),
	)

	repo := paymentRepository.NewPaymentRepository()
	paymentService := paymentService.NewService(repo)
	paymentApi := paymentApi.NewAPI(paymentService)

	paymentV1.RegisterPaymentServiceServer(s, paymentApi)

	// Debugging purposes
	reflection.Register(s)

	go func() {
		log.Printf("🚀 Payment gRPC server listening on port %d\n", grpcPort)
		err = s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve Payment gRPC: %v\n", err)
			return
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down Payment servers...")

	// Shutdown gRPC server
	s.GracefulStop()
	log.Println("✅ Payment gRPC server stopped")
}

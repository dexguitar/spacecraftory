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
	"github.com/dexguitar/spacecraftory/payment/internal/config"
	"github.com/dexguitar/spacecraftory/payment/internal/interceptor"
	paymentRepository "github.com/dexguitar/spacecraftory/payment/internal/repository/payment"
	paymentService "github.com/dexguitar/spacecraftory/payment/internal/service/payment"
	paymentV1 "github.com/dexguitar/spacecraftory/shared/pkg/proto/payment/v1"
)

const configPath = "./deploy/compose/payment/.env"

func main() {
	err := config.Load(configPath)
	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}

	lis, err := net.Listen("tcp", config.AppConfig().PaymentGRPC.Address())
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

	// Enable reflection for debugging
	reflection.Register(s)

	repo := paymentRepository.NewPaymentRepository()
	paymentService := paymentService.NewService(repo)
	paymentApi := paymentApi.NewAPI(paymentService)

	paymentV1.RegisterPaymentServiceServer(s, paymentApi)

	go func() {
		log.Printf("🚀 Payment gRPC server listening on %s\n", config.AppConfig().PaymentGRPC.Address())
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
	s.GracefulStop()
	log.Println("✅ Payment gRPC server stopped")
}

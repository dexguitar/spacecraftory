package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"slices"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"github.com/dexguitar/spacecraftory/payment/internal/interceptor"
	paymentV1 "github.com/dexguitar/spacecraftory/shared/pkg/proto/payment/v1"
)

const (
	grpcPort = 50052
	httpPort = 8082
)

var validPaymentMethods = []paymentV1.PaymentMethod{paymentV1.PaymentMethod_PAYMENT_METHOD_CARD, paymentV1.PaymentMethod_PAYMENT_METHOD_SBP, paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD, paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY}

// paymentService implements gRPC service for working with payment
type paymentService struct {
	paymentV1.UnimplementedPaymentServiceServer
}

// PayOrder returns a transaction uuid
func (s *paymentService) PayOrder(ctx context.Context, req *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	if req.OrderUuid == "" || req.UserUuid == "" {
		return nil, status.Errorf(codes.InvalidArgument, "order_uuid and user_uuid are required")
	}

	if !slices.Contains(validPaymentMethods, req.PaymentMethod) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid payment_method")
	}

	transaction_uuid := uuid.NewString()

	return &paymentV1.PayOrderResponse{
		TransactionUuid: transaction_uuid,
	}, nil
}

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

	// Enable reflection for debugging
	reflection.Register(s)

	// Register our service with mock data
	service := &paymentService{}

	paymentV1.RegisterPaymentServiceServer(s, service)

	go func() {
		log.Printf("🚀 Payment gRPC server listening on port %d\n", grpcPort)
		err = s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve Payment gRPC: %v\n", err)
			return
		}
	}()

	// Launch HTTP server with gRPC gateway
	var gwServer *http.Server
	go func() {
		// Create context with cancel
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Create new serve mux for HTTP requests
		mux := runtime.NewServeMux()

		// Create new dial options for HTTP requests
		opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

		// Register payment service handler from endpoint
		err = paymentV1.RegisterPaymentServiceHandlerFromEndpoint(
			ctx,
			mux,
			fmt.Sprintf("localhost:%d", grpcPort),
			opts,
		)
		if err != nil {
			log.Printf("Failed to register gateway for payment service: %v\n", err)
			return
		}

		// Create new HTTP server mux
		httpMux := http.NewServeMux()

		// Mount gRPC-Gateway at /api/
		httpMux.Handle("/api/", mux)

		// Serve Swagger JSON
		httpMux.HandleFunc("/apidocs.swagger.json", func(w http.ResponseWriter, r *http.Request) {
			log.Printf("📄 Serving swagger JSON: %s", r.URL.Path)
			http.ServeFile(w, r, "shared/pkg/swagger/payment/v1/payment.swagger.json")
		})

		// Serve Swagger UI HTML
		httpMux.HandleFunc("/swagger-ui.html", func(w http.ResponseWriter, r *http.Request) {
			log.Printf("📚 Serving swagger UI: %s", r.URL.Path)
			http.ServeFile(w, r, "shared/pkg/swagger/swagger-ui.html")
		})

		// Redirect root to Swagger UI
		httpMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			log.Printf("🔄 Request to: %s", r.URL.Path)
			if r.URL.Path == "/" {
				http.Redirect(w, r, "/swagger-ui.html", http.StatusMovedPermanently)
				return
			}
			http.NotFound(w, r)
		})

		// Create new HTTP server
		gwServer = &http.Server{
			Addr:              fmt.Sprintf(":%d", httpPort),
			Handler:           httpMux,
			ReadHeaderTimeout: 10 * time.Second,
		}

		// Start HTTP server
		log.Printf("🚀 Payment HTTP server with gRPC-Gateway listening on port %d\n", httpPort)
		log.Printf("📚 Swagger UI available at: http://localhost:%d/swagger-ui.html\n", httpPort)
		log.Printf("📄 Swagger JSON available at: http://localhost:%d/apidocs.swagger.json\n", httpPort)
		err = gwServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("Failed to serve Payment HTTP: %v\n", err)
			return
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down Payment servers...")

	// Shutdown HTTP server
	if gwServer != nil {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := gwServer.Shutdown(shutdownCtx); err != nil {
			log.Printf("Payment HTTP server shutdown error: %v", err)
		}
		log.Println("✅ Payment HTTP server stopped")
	}

	// Shutdown gRPC server
	s.GracefulStop()
	log.Println("✅ Payment gRPC server stopped")
}

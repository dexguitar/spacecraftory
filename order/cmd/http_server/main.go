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
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	customMiddleware "github.com/dexguitar/spacecraftory/order/internal/middleware"
	orderV1 "github.com/dexguitar/spacecraftory/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/dexguitar/spacecraftory/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/dexguitar/spacecraftory/shared/pkg/proto/payment/v1"
)

const (
	httpPort             = "8080"
	inventoryServiceAddr = "50051"
	paymentServiceAddr   = "50052"
	// HTTP Server timeouts
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

// OrderStorage represents a thread-safe storage for order data
type OrderStorage struct {
	mu     sync.RWMutex
	orders map[string]*orderV1.OrderDto
}

// NewOrderStorage creates a new storage for order data
func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		orders: make(map[string]*orderV1.OrderDto),
	}
}

// GetOrder returns the order by UUID
func (s *OrderStorage) GetOrder(orderUUID string) *orderV1.OrderDto {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order, ok := s.orders[orderUUID]
	if !ok {
		return nil
	}

	return order
}

// CreateOrder creates a new order
func (s *OrderStorage) CreateOrder(order *orderV1.OrderDto) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.orders[order.OrderUUID.String()] = order
}

// UpdateOrder updates an existing order
func (s *OrderStorage) UpdateOrder(order *orderV1.OrderDto) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.orders[order.OrderUUID.String()] = order
}

// OrderHandler implements the orderV1.Handler interface for handling requests to the Order API
type OrderHandler struct {
	storage         *OrderStorage
	inventoryClient inventoryV1.InventoryServiceClient
	paymentClient   paymentV1.PaymentServiceClient
}

// NewOrderHandler creates a new order handler
func NewOrderHandler(storage *OrderStorage, inventoryClient inventoryV1.InventoryServiceClient, paymentClient paymentV1.PaymentServiceClient) *OrderHandler {
	return &OrderHandler{
		storage:         storage,
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
	}
}

// CreateOrder handles the request to create a new order
func (h *OrderHandler) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	// Calculate total price by fetching part prices from inventory service
	var totalPrice float64

	for _, partUUID := range req.PartUuids {
		// Call inventory service to get part details
		resp, err := h.inventoryClient.GetPart(ctx, &inventoryV1.GetPartRequest{
			Uuid: partUUID.String(),
		})
		if err != nil {
			log.Printf("Error fetching part %s from inventory: %v", partUUID.String(), err)
			return &orderV1.BadRequestError{
				Code:    http.StatusBadRequest,
				Message: fmt.Sprintf("Failed to fetch part %s from inventory: %v", partUUID.String(), err),
			}, nil
		}

		// Add part price to total
		totalPrice += resp.Part.Price
	}

	// Generate new order UUID
	orderUUID := uuid.New()

	// Create the order object
	order := &orderV1.OrderDto{
		OrderUUID:  orderUUID,
		UserUUID:   req.UserUUID,
		PartUuids:  req.PartUuids,
		TotalPrice: totalPrice,
		Status:     orderV1.OrderStatusPENDINGPAYMENT,
	}

	// Save to storage
	h.storage.CreateOrder(order)

	// Return response
	return &orderV1.CreateOrderResponse{
		OrderUUID:  orderUUID,
		TotalPrice: totalPrice,
	}, nil
}

// GetOrderByUUID handles the request to get an order by UUID
func (h *OrderHandler) GetOrderByUUID(ctx context.Context, params orderV1.GetOrderByUUIDParams) (orderV1.GetOrderByUUIDRes, error) {
	order := h.storage.GetOrder(params.OrderUUID.String())
	if order == nil {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: fmt.Sprintf("Order %s not found", params.OrderUUID.String()),
		}, nil
	}

	return order, nil
}

// PayOrder handles the request to pay for an order
func (h *OrderHandler) PayOrder(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	order := h.storage.GetOrder(params.OrderUUID.String())
	if order == nil {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: fmt.Sprintf("Order %s not found", params.OrderUUID.String()),
		}, nil
	}

	if order.Status != orderV1.OrderStatusPENDINGPAYMENT {
		return &orderV1.ConflictError{
			Code:    409,
			Message: fmt.Sprintf("Order %s is not in pending payment status", params.OrderUUID.String()),
		}, nil
	}

	// Generate transaction UUID
	transactionUUID, err := h.paymentClient.PayOrder(ctx, &paymentV1.PayOrderRequest{
		OrderUuid:     order.OrderUUID.String(),
		UserUuid:      order.UserUUID.String(),
		PaymentMethod: paymentMethodToProto(req.PaymentMethod),
	})
	if err != nil {
		return &orderV1.InternalServerError{
			Code:    500,
			Message: fmt.Sprintf("Failed to pay order %s: %v", params.OrderUUID.String(), err),
		}, nil
	}

	// Update order with payment information
	order.Status = orderV1.OrderStatusPAID
	order.TransactionUUID = orderV1.NewOptNilUUID(uuid.MustParse(transactionUUID.TransactionUuid))
	order.PaymentMethod = orderV1.NewOptNilPaymentMethod(req.PaymentMethod)

	h.storage.UpdateOrder(order)

	return &orderV1.PayOrderResponse{
		TransactionUUID: uuid.MustParse(transactionUUID.TransactionUuid),
	}, nil
}

// CancelOrder handles the request to cancel an order
func (h *OrderHandler) CancelOrder(ctx context.Context, params orderV1.CancelOrderParams) (orderV1.CancelOrderRes, error) {
	order := h.storage.GetOrder(params.OrderUUID.String())
	if order == nil {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: fmt.Sprintf("Order %s not found", params.OrderUUID.String()),
		}, nil
	}

	// Check if order is already paid
	if order.Status == orderV1.OrderStatusPAID || order.Status == orderV1.OrderStatusASSEMBLED {
		return &orderV1.ConflictError{
			Code:    409,
			Message: "Order is already paid and cannot be cancelled",
		}, nil
	}

	// Update order status
	order.Status = orderV1.OrderStatusCANCELLED
	h.storage.UpdateOrder(order)

	return &orderV1.CancelOrderNoContent{}, nil
}

// NewError creates a new error in GenericError format
func (h *OrderHandler) NewError(ctx context.Context, err error) *orderV1.GenericErrorStatusCode {
	return &orderV1.GenericErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: orderV1.GenericError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		},
	}
}

func paymentMethodToProto(method orderV1.PaymentMethod) paymentV1.PaymentMethod {
	switch method {
	case orderV1.PaymentMethodCARD:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_CARD
	case orderV1.PaymentMethodSBP:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_SBP
	case orderV1.PaymentMethodCREDITCARD:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case orderV1.PaymentMethodINVESTORMONEY:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	default:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_UNKNOWN_UNSPECIFIED
	}
}

func main() {
	// Create gRPC connection to inventory service
	inventoryConn, err := grpc.NewClient(
		net.JoinHostPort("localhost", inventoryServiceAddr),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Failed to connect to inventory service: %v", err)
	}
	defer func() {
		if cerr := inventoryConn.Close(); cerr != nil {
			log.Printf("Failed to close inventory connection: %v", cerr)
		}
	}()

	// Create gRPC connection to inventory service
	paymentConn, err := grpc.NewClient(
		net.JoinHostPort("localhost", paymentServiceAddr),
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

	// Create inventory service client
	inventoryClient := inventoryV1.NewInventoryServiceClient(inventoryConn)
	log.Printf("✅ Connected to inventory service at %s", inventoryServiceAddr)

	paymentClient := paymentV1.NewPaymentServiceClient(paymentConn)
	log.Printf("✅ Connected to payment service at %s", paymentServiceAddr)

	// Create storage for order data
	storage := NewOrderStorage()

	// Create order API handler
	orderHandler := NewOrderHandler(storage, inventoryClient, paymentClient)

	// Create OpenAPI server
	orderServer, err := orderV1.NewServer(orderHandler)
	if err != nil {
		log.Printf("Error creating OpenAPI server: %v", err)
		return
	}

	// Initialize Chi router
	r := chi.NewRouter()

	// Add middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))
	r.Use(customMiddleware.RequestLogger)

	// Mount OpenAPI handlers
	r.Mount("/", orderServer)

	// Start HTTP server
	server := &http.Server{
		Addr:              net.JoinHostPort("localhost", httpPort),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout, // Protection from Slowloris attacks - a type of DDoS attack where
		// the attacker intentionally sends HTTP headers slowly, keeping connections open and exhausting
		// the server's connection pool. ReadHeaderTimeout forcibly closes the connection
		// if the client fails to send all headers within the allotted time.
	}

	// Start server in a separate goroutine
	go func() {
		log.Printf("🚀 HTTP server started on port %s\n", httpPort)
		err = server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("❌ Server startup error: %v\n", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")

	// Create context with timeout for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("❌ Error during server shutdown: %v\n", err)
	}

	log.Println("✅ Server stopped")
}

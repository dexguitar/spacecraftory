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

var (
	ErrOrderNotFound      = errors.New("order not found")
	ErrOrderAlreadyExists = errors.New("order already exists")
)

// OrderStorage represents a thread-safe storage for order data
type OrderStorage struct {
	mu     sync.RWMutex
	orders map[string]*orderV1.OrderDto
}

func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		orders: make(map[string]*orderV1.OrderDto),
	}
}

func (s *OrderStorage) GetOrder(orderUUID string) (*orderV1.OrderDto, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order, ok := s.orders[orderUUID]
	if !ok {
		return &orderV1.OrderDto{}, ErrOrderNotFound
	}

	return order, nil
}

func (s *OrderStorage) CreateOrder(order *orderV1.OrderDto) (*orderV1.OrderDto, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.orders[order.OrderUUID.String()] = order

	return order, nil
}

func (s *OrderStorage) UpdateOrder(order *orderV1.OrderDto) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.orders[order.OrderUUID.String()]; !ok {
		return ErrOrderNotFound
	}

	return nil
}

// OrderHandler implements the orderV1.Handler interface for handling requests to the Order API
type OrderHandler struct {
	storage         *OrderStorage
	inventoryClient inventoryV1.InventoryServiceClient
	paymentClient   paymentV1.PaymentServiceClient
}

func NewOrderHandler(storage *OrderStorage, inventoryClient inventoryV1.InventoryServiceClient, paymentClient paymentV1.PaymentServiceClient) *OrderHandler {
	return &OrderHandler{
		storage:         storage,
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
	}
}

func (h *OrderHandler) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	resp, err := h.inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
		Filter: &inventoryV1.PartsFilter{
			Uuids: convertPartUUIDsToProto(req.PartUuids),
		},
	})
	if err != nil {
		return &orderV1.BadRequestError{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf("Failed to fetch parts from inventory: %v", err),
		}, nil
	}

	if len(resp.Parts) != len(req.PartUuids) {
		return &orderV1.BadRequestError{
			Code:    http.StatusBadRequest,
			Message: "Some parts were not found",
		}, nil
	}

	orderTotalPrice := countTotalPartsPrice(resp.Parts)

	// Generate new order UUID
	orderUUID := uuid.New()

	order := &orderV1.OrderDto{
		OrderUUID:  orderUUID,
		UserUUID:   req.UserUUID,
		PartUuids:  req.PartUuids,
		TotalPrice: orderTotalPrice,
		Status:     orderV1.OrderStatusPENDINGPAYMENT,
	}

	_, err = h.storage.CreateOrder(order)
	if err != nil {
		return &orderV1.InternalServerError{
			Code:    500,
			Message: fmt.Sprintf("Failed to create order: %v", err),
		}, nil
	}

	return &orderV1.CreateOrderResponse{
		OrderUUID:  orderUUID,
		TotalPrice: orderTotalPrice,
	}, nil
}

func (h *OrderHandler) GetOrderByUUID(ctx context.Context, params orderV1.GetOrderByUUIDParams) (orderV1.GetOrderByUUIDRes, error) {
	order, err := h.storage.GetOrder(params.OrderUUID.String())
	if err != nil {
		if errors.Is(err, ErrOrderNotFound) {
			return &orderV1.NotFoundError{
				Code:    404,
				Message: fmt.Sprintf("Order %s not found", params.OrderUUID.String()),
			}, nil
		} else {
			return &orderV1.InternalServerError{
				Code:    500,
				Message: fmt.Sprintf("Failed to get order %s: %v", params.OrderUUID.String(), err),
			}, nil
		}
	}

	return order, nil
}

func (h *OrderHandler) PayOrder(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	order, err := h.storage.GetOrder(params.OrderUUID.String())
	if err != nil {
		if errors.Is(err, ErrOrderNotFound) {
			return &orderV1.NotFoundError{
				Code:    404,
				Message: fmt.Sprintf("Order %s not found", params.OrderUUID.String()),
			}, nil
		} else {
			return &orderV1.InternalServerError{
				Code:    500,
				Message: fmt.Sprintf("Failed to get order %s: %v", params.OrderUUID.String(), err),
			}, nil
		}
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

	order.Status = orderV1.OrderStatusPAID
	parsedUUID, err := uuid.Parse(transactionUUID.TransactionUuid)
	if err != nil {
		return &orderV1.InternalServerError{
			Code:    500,
			Message: fmt.Sprintf("Failed to parse transaction UUID %s: %v", transactionUUID.TransactionUuid, err),
		}, nil
	}
	order.TransactionUUID = orderV1.NewOptNilUUID(parsedUUID)
	order.PaymentMethod = orderV1.NewOptNilPaymentMethod(req.PaymentMethod)

	err = h.storage.UpdateOrder(order)
	if err != nil {
		return &orderV1.InternalServerError{
			Code:    500,
			Message: fmt.Sprintf("Failed to update order %s: %v", params.OrderUUID.String(), err),
		}, nil
	}

	return &orderV1.PayOrderResponse{
		TransactionUUID: uuid.MustParse(transactionUUID.TransactionUuid),
	}, nil
}

func (h *OrderHandler) CancelOrder(ctx context.Context, params orderV1.CancelOrderParams) (orderV1.CancelOrderRes, error) {
	order, err := h.storage.GetOrder(params.OrderUUID.String())
	if err != nil {
		if errors.Is(err, ErrOrderNotFound) {
			return &orderV1.NotFoundError{
				Code:    404,
				Message: fmt.Sprintf("Order %s not found", params.OrderUUID.String()),
			}, nil
		} else {
			return &orderV1.InternalServerError{
				Code:    500,
				Message: fmt.Sprintf("Failed to get order %s: %v", params.OrderUUID.String(), err),
			}, nil
		}
	}

	if order.Status == orderV1.OrderStatusPAID || order.Status == orderV1.OrderStatusCANCELLED {
		return &orderV1.ConflictError{
			Code:    409,
			Message: "Order has already been paid or cancelled",
		}, nil
	}

	order.Status = orderV1.OrderStatusCANCELLED
	err = h.storage.UpdateOrder(order)
	if err != nil {
		return &orderV1.InternalServerError{
			Code:    500,
			Message: fmt.Sprintf("Failed to update order %s: %v", params.OrderUUID.String(), err),
		}, nil
	}

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

// CreateOrder helper functions
func convertPartUUIDsToProto(partUUIDs []uuid.UUID) []string {
	result := make([]string, 0, len(partUUIDs))
	for _, partUUID := range partUUIDs {
		result = append(result, partUUID.String())
	}
	return result
}

func countTotalPartsPrice(parts []*inventoryV1.Part) float64 {
	var totalPrice float64
	for _, part := range parts {
		totalPrice += part.Price
	}
	return totalPrice
}

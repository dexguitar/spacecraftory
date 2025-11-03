package app

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderV1API "github.com/dexguitar/spacecraftory/order/internal/api/order/v1"
	inventoryClientImpl "github.com/dexguitar/spacecraftory/order/internal/client/grpc/inventory/v1"
	paymentClientImpl "github.com/dexguitar/spacecraftory/order/internal/client/grpc/payment/v1"
	invClient "github.com/dexguitar/spacecraftory/order/internal/client/inventory"
	payClient "github.com/dexguitar/spacecraftory/order/internal/client/payment"
	"github.com/dexguitar/spacecraftory/order/internal/config"
	"github.com/dexguitar/spacecraftory/order/internal/repository"
	orderRepository "github.com/dexguitar/spacecraftory/order/internal/repository/order"
	"github.com/dexguitar/spacecraftory/order/internal/service"
	orderService "github.com/dexguitar/spacecraftory/order/internal/service/order"
	"github.com/dexguitar/spacecraftory/platform/pkg/closer"
	orderV1 "github.com/dexguitar/spacecraftory/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/dexguitar/spacecraftory/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/dexguitar/spacecraftory/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	orderV1API orderV1.Handler

	orderService service.OrderService

	orderRepository repository.OrderRepository

	inventoryClient invClient.InventoryClient
	paymentClient   payClient.PaymentClient

	inventoryGRPCConn *grpc.ClientConn
	paymentGRPCConn   *grpc.ClientConn

	pgPool *pgxpool.Pool
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) OrderV1API(ctx context.Context) orderV1.Handler {
	if d.orderV1API == nil {
		d.orderV1API = orderV1API.NewAPI(d.OrderService(ctx))
	}

	return d.orderV1API
}

func (d *diContainer) OrderService(ctx context.Context) service.OrderService {
	if d.orderService == nil {
		d.orderService = orderService.NewService(
			d.OrderRepository(ctx),
			d.InventoryClient(ctx),
			d.PaymentClient(ctx),
		)
	}

	return d.orderService
}

func (d *diContainer) OrderRepository(ctx context.Context) repository.OrderRepository {
	if d.orderRepository == nil {
		d.orderRepository = orderRepository.NewOrderRepository(d.PgPool(ctx))
	}

	return d.orderRepository
}

func (d *diContainer) InventoryClient(ctx context.Context) invClient.InventoryClient {
	if d.inventoryClient == nil {
		grpcClient := inventoryV1.NewInventoryServiceClient(d.InventoryGRPCConn(ctx))
		d.inventoryClient = inventoryClientImpl.NewInventoryClient(grpcClient)
	}

	return d.inventoryClient
}

func (d *diContainer) PaymentClient(ctx context.Context) payClient.PaymentClient {
	if d.paymentClient == nil {
		grpcClient := paymentV1.NewPaymentServiceClient(d.PaymentGRPCConn(ctx))
		d.paymentClient = paymentClientImpl.NewPaymentClient(grpcClient)
	}

	return d.paymentClient
}

func (d *diContainer) InventoryGRPCConn(_ context.Context) *grpc.ClientConn {
	if d.inventoryGRPCConn == nil {
		conn, err := grpc.NewClient(
			fmt.Sprintf(":%s", config.AppConfig().GRPCClient.InventoryAddress()),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to connect to inventory service: %s", err.Error()))
		}

		closer.AddNamed("Inventory gRPC connection", func(ctx context.Context) error {
			return conn.Close()
		})

		d.inventoryGRPCConn = conn
	}

	return d.inventoryGRPCConn
}

func (d *diContainer) PaymentGRPCConn(_ context.Context) *grpc.ClientConn {
	if d.paymentGRPCConn == nil {
		conn, err := grpc.NewClient(
			fmt.Sprintf(":%s", config.AppConfig().GRPCClient.PaymentAddress()),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to connect to payment service: %s", err.Error()))
		}

		closer.AddNamed("Payment gRPC connection", func(ctx context.Context) error {
			return conn.Close()
		})

		d.paymentGRPCConn = conn
	}

	return d.paymentGRPCConn
}

func (d *diContainer) PgPool(ctx context.Context) *pgxpool.Pool {
	if d.pgPool == nil {
		dbURI := config.AppConfig().Postgres.Address()

		pool, err := pgxpool.New(ctx, dbURI)
		if err != nil {
			panic(fmt.Sprintf("failed to create connection pool: %s", err.Error()))
		}

		closer.AddNamed("PostgreSQL connection pool", func(ctx context.Context) error {
			pool.Close()
			return nil
		})

		d.pgPool = pool
	}

	return d.pgPool
}

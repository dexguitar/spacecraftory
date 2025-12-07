package app

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	inventoryV1API "github.com/dexguitar/spacecraftory/inventory/internal/api/inventory/v1"
	iamClient "github.com/dexguitar/spacecraftory/inventory/internal/client/grpc/iam/v1"
	"github.com/dexguitar/spacecraftory/inventory/internal/config"
	"github.com/dexguitar/spacecraftory/inventory/internal/repository"
	inventoryRepository "github.com/dexguitar/spacecraftory/inventory/internal/repository/inventory"
	"github.com/dexguitar/spacecraftory/inventory/internal/service"
	inventoryService "github.com/dexguitar/spacecraftory/inventory/internal/service/inventory"
	"github.com/dexguitar/spacecraftory/platform/pkg/closer"
	authV1 "github.com/dexguitar/spacecraftory/shared/pkg/proto/auth/v1"
	inventoryV1 "github.com/dexguitar/spacecraftory/shared/pkg/proto/inventory/v1"
)

type diContainer struct {
	inventoryV1API inventoryV1.InventoryServiceServer
	iamClient      iamClient.IAMClient
	iamGRPCClient  authV1.AuthServiceClient
	iamGRPCConn    *grpc.ClientConn

	inventoryService service.InventoryService

	inventoryRepository repository.InventoryRepository

	mongoDBClient *mongo.Client
	mongoDBHandle *mongo.Database
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) InventoryV1API(ctx context.Context) inventoryV1.InventoryServiceServer {
	if d.inventoryV1API == nil {
		d.inventoryV1API = inventoryV1API.NewAPI(d.InventoryService(ctx))
	}

	return d.inventoryV1API
}

func (d *diContainer) IAMClient(ctx context.Context) iamClient.IAMClient {
	if d.iamClient == nil {
		d.iamClient = iamClient.NewIAMClient(d.IAMGRPCClient(ctx))
	}

	return d.iamClient
}

func (d *diContainer) IAMGRPCClient(ctx context.Context) authV1.AuthServiceClient {
	if d.iamGRPCClient == nil {
		d.iamGRPCClient = authV1.NewAuthServiceClient(d.IAMGRPCConn(ctx))
	}

	return d.iamGRPCClient
}

func (d *diContainer) IAMGRPCConn(ctx context.Context) *grpc.ClientConn {
	if d.iamGRPCConn == nil {
		conn, err := grpc.NewClient(
			config.AppConfig().IAMClientGRPC.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to connect to IAM service: %s", err.Error()))
		}

		closer.AddNamed("IAM gRPC connection", func(ctx context.Context) error {
			return conn.Close()
		})

		d.iamGRPCConn = conn
	}

	return d.iamGRPCConn
}

func (d *diContainer) InventoryService(ctx context.Context) service.InventoryService {
	if d.inventoryService == nil {
		d.inventoryService = inventoryService.NewService(d.InventoryRepository(ctx))
	}

	return d.inventoryService
}

func (d *diContainer) InventoryRepository(ctx context.Context) repository.InventoryRepository {
	if d.inventoryRepository == nil {
		d.inventoryRepository = inventoryRepository.NewInventoryRepository(ctx, d.MongoDBHandle(ctx))
	}

	return d.inventoryRepository
}

func (d *diContainer) MongoDBClient(ctx context.Context) *mongo.Client {
	if d.mongoDBClient == nil {
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.AppConfig().Mongo.URI()))
		if err != nil {
			panic(fmt.Sprintf("failed to connect to MongoDB: %s\n", err.Error()))
		}

		err = client.Ping(ctx, readpref.Primary())
		if err != nil {
			panic(fmt.Sprintf("failed to ping MongoDB: %v\n", err))
		}

		closer.AddNamed("MongoDB client", func(ctx context.Context) error {
			return client.Disconnect(ctx)
		})

		d.mongoDBClient = client
	}

	return d.mongoDBClient
}

func (d *diContainer) MongoDBHandle(ctx context.Context) *mongo.Database {
	if d.mongoDBHandle == nil {
		d.mongoDBHandle = d.MongoDBClient(ctx).Database(config.AppConfig().Mongo.DatabaseName())
	}

	return d.mongoDBHandle
}

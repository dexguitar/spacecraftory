package integration

import (
	"context"
	"os"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"

	"github.com/dexguitar/spacecraftory/platform/pkg/logger"
	"github.com/dexguitar/spacecraftory/platform/pkg/testcontainers"
	"github.com/dexguitar/spacecraftory/platform/pkg/testcontainers/app"
	"github.com/dexguitar/spacecraftory/platform/pkg/testcontainers/mongo"
	"github.com/dexguitar/spacecraftory/platform/pkg/testcontainers/network"
	"github.com/dexguitar/spacecraftory/platform/pkg/testcontainers/path"
)

const (
	// Parameters for containers
	inventoryAppName    = "inventory-app"
	inventoryDockerfile = "deploy/docker/inventory/Dockerfile"

	// Application environment variables
	grpcPortKey = "GRPC_PORT"

	// Environment variable values
	loggerLevelValue = "debug"
	startupTimeout   = 3 * time.Minute
)

// TestEnvironment is a structure for storing test environment resources
type TestEnvironment struct {
	Network *network.Network
	Mongo   *mongo.Container
	App     *app.Container
}

// setupTestEnvironment prepares the test environment: network, containers and returns a structure with resources
func setupTestEnvironment(ctx context.Context) *TestEnvironment {
	logger.Info(ctx, "🚀 Preparing test environment...")

	// Step 1: Create shared Docker network
	generatedNetwork, err := network.NewNetwork(ctx, projectName)
	if err != nil {
		logger.Fatal(ctx, "failed to create shared network", zap.Error(err))
	}
	logger.Info(ctx, "✅ Network successfully created")

	// Get environment variables for MongoDB with existence check
	mongoUsername := getEnvWithLogging(ctx, testcontainers.MongoUsernameKey)
	mongoPassword := getEnvWithLogging(ctx, testcontainers.MongoPasswordKey)
	mongoImageName := getEnvWithLogging(ctx, testcontainers.MongoImageNameKey)
	mongoDatabase := getEnvWithLogging(ctx, testcontainers.MongoDatabaseKey)

	// Get gRPC port for waitStrategy
	grpcPort := getEnvWithLogging(ctx, grpcPortKey)

	// Step 2: Start MongoDB container
	generatedMongo, err := mongo.NewContainer(ctx,
		mongo.WithNetworkName(generatedNetwork.Name()),
		mongo.WithContainerName(testcontainers.MongoContainerName),
		mongo.WithImageName(mongoImageName),
		mongo.WithDatabase(mongoDatabase),
		mongo.WithAuth(mongoUsername, mongoPassword),
		mongo.WithLogger(logger.Logger()),
	)
	if err != nil {
		cleanupTestEnvironment(ctx, &TestEnvironment{Network: generatedNetwork})
		logger.Fatal(ctx, "failed to start MongoDB container", zap.Error(err))
	}
	logger.Info(ctx, "✅ MongoDB container successfully started")

	// Step 3: Start application container
	projectRoot := path.GetProjectRoot()

	appEnv := map[string]string{
		// Override MongoDB host to connect to container from testcontainers
		testcontainers.MongoHostKey: generatedMongo.Config().ContainerName,
	}

	// Create custom wait strategy with increased timeout
	waitStrategy := wait.ForListeningPort(nat.Port(grpcPort + "/tcp")).
		WithStartupTimeout(startupTimeout)

	appContainer, err := app.NewContainer(ctx,
		app.WithName(inventoryAppName),
		app.WithPort(grpcPort),
		app.WithDockerfile(projectRoot, inventoryDockerfile),
		app.WithNetwork(generatedNetwork.Name()),
		app.WithEnv(appEnv),
		app.WithLogOutput(os.Stdout),
		app.WithStartupWait(waitStrategy),
		app.WithLogger(logger.Logger()),
	)
	if err != nil {
		cleanupTestEnvironment(ctx, &TestEnvironment{Network: generatedNetwork, Mongo: generatedMongo})
		logger.Fatal(ctx, "failed to start application container", zap.Error(err))
	}
	logger.Info(ctx, "✅ Application container successfully started")

	logger.Info(ctx, "🎉 Test environment is ready")
	return &TestEnvironment{
		Network: generatedNetwork,
		Mongo:   generatedMongo,
		App:     appContainer,
	}
}

// getEnvWithLogging returns environment variable value with logging
func getEnvWithLogging(ctx context.Context, key string) string {
	value := os.Getenv(key)
	if value == "" {
		logger.Warn(ctx, "Environment variable is not set", zap.String("key", key))
	}

	return value
}

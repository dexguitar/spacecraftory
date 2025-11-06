package integration

import (
	"context"

	"go.uber.org/zap"

	"github.com/dexguitar/spacecraftory/platform/pkg/logger"
)

// teardownTestEnvironment releases all test environment resources
func teardownTestEnvironment(ctx context.Context, env *TestEnvironment) {
	log := logger.Logger()
	log.Info(ctx, "🧹 Cleaning up test environment...")

	cleanupTestEnvironment(ctx, env)

	log.Info(ctx, "✅ Test environment successfully cleaned up")
}

// cleanupTestEnvironment is a helper function for releasing resources
func cleanupTestEnvironment(ctx context.Context, env *TestEnvironment) {
	if env.App != nil {
		if err := env.App.Terminate(ctx); err != nil {
			logger.Error(ctx, "failed to stop application container", zap.Error(err))
		} else {
			logger.Info(ctx, "🛑 Application container stopped")
		}
	}

	if env.Mongo != nil {
		if err := env.Mongo.Terminate(ctx); err != nil {
			logger.Error(ctx, "failed to stop MongoDB container", zap.Error(err))
		} else {
			logger.Info(ctx, "🛑 MongoDB container stopped")
		}
	}

	if env.Network != nil {
		if err := env.Network.Remove(ctx); err != nil {
			logger.Error(ctx, "failed to remove network", zap.Error(err))
		} else {
			logger.Info(ctx, "🛑 Network removed")
		}
	}
}

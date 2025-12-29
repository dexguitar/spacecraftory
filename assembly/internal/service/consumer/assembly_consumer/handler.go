package assembly_consumer

import (
	"context"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/dexguitar/spacecraftory/assembly/internal/metrics"
	"github.com/dexguitar/spacecraftory/assembly/internal/model"
	"github.com/dexguitar/spacecraftory/platform/pkg/kafka"
	"github.com/dexguitar/spacecraftory/platform/pkg/logger"
)

func (s *orderPaidConsumerService) OrderPaidHandler(ctx context.Context, msg kafka.Message) error {
	event, err := s.orderPaidDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode OrderPaid", zap.Error(err))
		return err
	}

	logger.Info(ctx, "Processing message",
		zap.String("topic", msg.Topic),
		zap.Any("partition", msg.Partition),
		zap.Any("offset", msg.Offset),
		zap.String("event_uuid", event.EventUUID),
		zap.String("order_uuid", event.OrderUUID),
		zap.String("user_uuid", event.UserUUID),
		zap.String("payment_method", event.PaymentMethod),
		zap.String("transaction_uuid", event.TransactionUUID),
	)

	// Generate random assembly duration between 5 and 10 seconds
	//nolint:gosec
	randomDuration := time.Duration(5+rand.Intn(6)) * time.Second

	// Start measuring assembly time
	assemblyStart := time.Now()

	// Simulate assembly process with random duration
	//nolint:forbidigo
	time.Sleep(randomDuration)

	// Record assembly duration metric
	assemblyTime := time.Since(assemblyStart)
	if metrics.AssemblyDuration != nil {
		metrics.AssemblyDuration.Record(ctx, assemblyTime.Seconds())
	}

	logger.Info(ctx, "Assembly completed",
		zap.String("order_uuid", event.OrderUUID),
		zap.Duration("assembly_time", assemblyTime))

	shipAssembledEvent := model.ShipAssembledEvent{
		EventUUID:    uuid.New().String(),
		OrderUUID:    event.OrderUUID,
		UserUUID:     event.UserUUID,
		BuildTimeSec: int64(assemblyTime.Seconds()),
	}

	if err := s.producerService.ProduceShipAssembled(ctx, shipAssembledEvent); err != nil {
		logger.Error(ctx, "Failed to produce ShipAssembled event",
			zap.String("order_uuid", event.OrderUUID),
			zap.Error(err))
		return err
	}

	logger.Info(ctx, "ShipAssembled event produced successfully",
		zap.String("event_uuid", shipAssembledEvent.EventUUID),
		zap.String("order_uuid", shipAssembledEvent.OrderUUID))

	return nil
}

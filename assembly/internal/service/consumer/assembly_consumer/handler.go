package assembly_consumer

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/dexguitar/spacecraftory/assembly/internal/model"
	"github.com/dexguitar/spacecraftory/platform/pkg/kafka"
	"github.com/dexguitar/spacecraftory/platform/pkg/logger"
)

const assemblyDuration = 10 * time.Second

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

	// Simulate assembly process
	//nolint:forbidigo
	time.Sleep(assemblyDuration)

	logger.Info(ctx, "Assembly completed", zap.String("order_uuid", event.OrderUUID))

	shipAssembledEvent := model.ShipAssembledEvent{
		EventUUID:    uuid.New().String(),
		OrderUUID:    event.OrderUUID,
		UserUUID:     event.UserUUID,
		BuildTimeSec: int64(assemblyDuration.Seconds()),
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

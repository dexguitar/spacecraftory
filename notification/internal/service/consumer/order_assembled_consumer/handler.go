package order_assembled_consumer

import (
	"context"

	"go.uber.org/zap"

	wrappedKafka "github.com/dexguitar/spacecraftory/platform/pkg/kafka"
	"github.com/dexguitar/spacecraftory/platform/pkg/logger"
)

func (s *orderAssembledConsumerService) OrderAssembledHandler(ctx context.Context, msg wrappedKafka.Message) error {
	event, err := s.orderAssembledDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode OrderAssembled", zap.Error(err))
		return err
	}

	err = s.telegramService.SendOrderAssembledNotification(ctx, event)
	if err != nil {
		logger.Error(ctx, "Failed to send order assembled notification", zap.Error(err))
		return err
	}

	logger.Info(ctx, "Processing message",
		zap.String("topic", msg.Topic),
		zap.Any("partition", msg.Partition),
		zap.Any("offset", msg.Offset),
		zap.String("event_uuid", event.EventUUID),
		zap.String("order_uuid", event.OrderUUID),
		zap.String("user_uuid", event.UserUUID),
		zap.Int64("build_time_sec", event.BuildTimeSec),
	)

	return nil
}

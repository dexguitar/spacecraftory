package order_paid_consumer

import (
	"context"

	"go.uber.org/zap"

	wrappedKafka "github.com/dexguitar/spacecraftory/platform/pkg/kafka"
	"github.com/dexguitar/spacecraftory/platform/pkg/logger"
)

func (s *orderPaidConsumerService) OrderPaidHandler(ctx context.Context, msg wrappedKafka.Message) error {
	event, err := s.orderPaidDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode OrderPaid", zap.Error(err))
		return err
	}

	err = s.telegramService.SendOrderPaidNotification(ctx, event)
	if err != nil {
		logger.Error(ctx, "Failed to send order paid notification", zap.Error(err))
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

	return nil
}

package order_consumer

import (
	"context"

	"go.uber.org/zap"

	"github.com/dexguitar/spacecraftory/order/internal/model"
	"github.com/dexguitar/spacecraftory/platform/pkg/kafka"
	"github.com/dexguitar/spacecraftory/platform/pkg/logger"
)

func (s *service) OrderHandler(ctx context.Context, msg kafka.Message) error {
	event, err := s.orderAssembledDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode OrderAssembled", zap.Error(err))
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

	order, err := s.orderRepository.GetOrder(ctx, event.OrderUUID)
	if err != nil {
		logger.Error(ctx, "Failed to get order", zap.Error(err))
		return err
	}

	order.OrderStatus = model.OrderStatusASSEMBLED

	err = s.orderRepository.UpdateOrder(ctx, order)
	if err != nil {
		logger.Error(ctx, "Failed to update order", zap.Error(err))
		return err
	}

	logger.Info(ctx, "Order updated successfully", zap.String("order_uuid", order.OrderUUID))

	return nil
}

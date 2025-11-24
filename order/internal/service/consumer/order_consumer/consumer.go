package order_consumer

import (
	"context"

	"go.uber.org/zap"

	kafkaConverter "github.com/dexguitar/spacecraftory/order/internal/converter/kafka"
	"github.com/dexguitar/spacecraftory/order/internal/repository"
	"github.com/dexguitar/spacecraftory/platform/pkg/kafka"
	"github.com/dexguitar/spacecraftory/platform/pkg/logger"
)

type service struct {
	orderAssembledConsumer kafka.Consumer
	orderAssembledDecoder  kafkaConverter.OrderAssembledDecoder
	orderRepository        repository.OrderRepository
}

func NewService(orderAssembledConsumer kafka.Consumer, orderAssembledDecoder kafkaConverter.OrderAssembledDecoder, orderRepository repository.OrderRepository) *service {
	return &service{
		orderAssembledConsumer: orderAssembledConsumer,
		orderAssembledDecoder:  orderAssembledDecoder,
		orderRepository:        orderRepository,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "🚀 OrderAssembled Kafka consumer running")

	err := s.orderAssembledConsumer.Consume(ctx, s.OrderHandler)
	if err != nil {
		logger.Error(ctx, "Consume from order.assembled topic error", zap.Error(err))
		return err
	}

	return nil
}

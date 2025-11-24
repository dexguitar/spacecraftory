package order_assembled_consumer

import (
	"context"

	"go.uber.org/zap"

	kafkaConverter "github.com/dexguitar/spacecraftory/notification/internal/converter/kafka"
	"github.com/dexguitar/spacecraftory/notification/internal/service"
	wrappedKafka "github.com/dexguitar/spacecraftory/platform/pkg/kafka"
	"github.com/dexguitar/spacecraftory/platform/pkg/logger"
)

type orderAssembledConsumerService struct {
	orderAssembledConsumer wrappedKafka.Consumer
	orderAssembledDecoder  kafkaConverter.OrderAssembledDecoder
	telegramService        service.TelegramService
}

func NewService(orderAssembledConsumer wrappedKafka.Consumer, orderAssembledDecoder kafkaConverter.OrderAssembledDecoder, telegramService service.TelegramService) *orderAssembledConsumerService {
	return &orderAssembledConsumerService{
		orderAssembledConsumer: orderAssembledConsumer,
		orderAssembledDecoder:  orderAssembledDecoder,
		telegramService:        telegramService,
	}
}

func (s *orderAssembledConsumerService) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "🚀 Order assembled Kafka consumer running")

	err := s.orderAssembledConsumer.Consume(ctx, s.OrderAssembledHandler)
	if err != nil {
		logger.Error(ctx, "Consume from order assembled topic error", zap.Error(err))
		return err
	}

	return nil
}

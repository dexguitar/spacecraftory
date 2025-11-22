package order_paid_consumer

import (
	"context"

	"go.uber.org/zap"

	kafkaConverter "github.com/dexguitar/spacecraftory/notification/internal/converter/kafka"
	"github.com/dexguitar/spacecraftory/notification/internal/service"
	wrappedKafka "github.com/dexguitar/spacecraftory/platform/pkg/kafka"
	"github.com/dexguitar/spacecraftory/platform/pkg/logger"
)

type orderPaidConsumerService struct {
	orderPaidConsumer wrappedKafka.Consumer
	orderPaidDecoder  kafkaConverter.OrderPaidDecoder
	telegramService   service.TelegramService
}

func NewService(orderPaidConsumer wrappedKafka.Consumer, orderPaidDecoder kafkaConverter.OrderPaidDecoder, telegramService service.TelegramService) *orderPaidConsumerService {
	return &orderPaidConsumerService{
		orderPaidConsumer: orderPaidConsumer,
		orderPaidDecoder:  orderPaidDecoder,
		telegramService:   telegramService,
	}
}

func (s *orderPaidConsumerService) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "🚀 Order assembled Kafka consumer running")

	err := s.orderPaidConsumer.Consume(ctx, s.OrderPaidHandler)
	if err != nil {
		logger.Error(ctx, "Consume from order paid topic error", zap.Error(err))
		return err
	}

	return nil
}

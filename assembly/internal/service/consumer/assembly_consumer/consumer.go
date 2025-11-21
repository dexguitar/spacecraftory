package assembly_consumer

import (
	"context"

	"go.uber.org/zap"

	kafkaConverter "github.com/dexguitar/spacecraftory/assembly/internal/converter/kafka"
	"github.com/dexguitar/spacecraftory/assembly/internal/service"
	"github.com/dexguitar/spacecraftory/platform/pkg/kafka"
	"github.com/dexguitar/spacecraftory/platform/pkg/logger"
)

type orderPaidConsumerService struct {
	orderPaidConsumer kafka.Consumer
	orderPaidDecoder  kafkaConverter.OrderPaidDecoder
	producerService   service.ProducerService
}

func NewService(
	orderPaidConsumer kafka.Consumer,
	orderPaidDecoder kafkaConverter.OrderPaidDecoder,
	producerService service.ProducerService,
) *orderPaidConsumerService {
	return &orderPaidConsumerService{
		orderPaidConsumer: orderPaidConsumer,
		orderPaidDecoder:  orderPaidDecoder,
		producerService:   producerService,
	}
}

func (s *orderPaidConsumerService) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting order paid consumer service")

	err := s.orderPaidConsumer.Consume(ctx, s.OrderPaidHandler)
	if err != nil {
		logger.Error(ctx, "Consume from order paid topic error", zap.Error(err))
		return err
	}

	return nil
}

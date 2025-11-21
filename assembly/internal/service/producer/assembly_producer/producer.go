package assembly_producer

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"github.com/dexguitar/spacecraftory/assembly/internal/model"
	"github.com/dexguitar/spacecraftory/platform/pkg/kafka"
	"github.com/dexguitar/spacecraftory/platform/pkg/logger"
	eventsV1 "github.com/dexguitar/spacecraftory/shared/pkg/proto/events/v1"
)

type shipAssembledProducerService struct {
	orderAssembledProducer kafka.Producer
}

func NewService(orderAssembledProducer kafka.Producer) *shipAssembledProducerService {
	return &shipAssembledProducerService{
		orderAssembledProducer: orderAssembledProducer,
	}
}

func (p *shipAssembledProducerService) ProduceShipAssembled(ctx context.Context, event model.ShipAssembledEvent) error {
	msg := &eventsV1.ShipAssembled{
		EventUuid:    event.EventUUID,
		OrderUuid:    event.OrderUUID,
		UserUuid:     event.UserUUID,
		BuildTimeSec: event.BuildTimeSec,
	}

	payload, err := proto.Marshal(msg)
	if err != nil {
		logger.Error(ctx, "failed to marshal ShipAssembled", zap.Error(err))
		return err
	}

	err = p.orderAssembledProducer.Send(ctx, []byte(event.EventUUID), payload)
	if err != nil {
		logger.Error(ctx, "failed to publish ShipAssembled", zap.Error(err))
		return err
	}

	return nil
}

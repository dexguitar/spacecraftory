package service

import (
	"context"

	"github.com/dexguitar/spacecraftory/assembly/internal/model"
)

type ProducerService interface {
	ProduceShipAssembled(ctx context.Context, event model.ShipAssembledEvent) error
}

type ConsumerService interface {
	RunConsumer(ctx context.Context) error
}

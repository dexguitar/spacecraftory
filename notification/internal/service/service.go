package service

import (
	"context"

	"github.com/dexguitar/spacecraftory/notification/internal/model"
)

type ConsumerService interface {
	RunConsumer(ctx context.Context) error
}

type TelegramService interface {
	SendOrderPaidNotification(ctx context.Context, event model.OrderPaidEvent) error
	SendOrderAssembledNotification(ctx context.Context, event model.OrderAssembledEvent) error
}

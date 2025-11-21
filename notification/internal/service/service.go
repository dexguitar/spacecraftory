package service

import (
	"context"
)

type ConsumerService interface {
	RunConsumer(ctx context.Context) error
}

type TelegramService interface {
	SendOrderPaidNotification(ctx context.Context, orderUUID, userUUID, paymentMethod, transactionUUID string) error
	SendOrderAssembledNotification(ctx context.Context, orderUUID, userUUID string, buildTimeSec int64) error
}

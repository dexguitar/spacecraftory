package service

import (
	"context"

	"github.com/dexguitar/spacecraftory/order/internal/model"
)

type OrderService interface {
	CreateOrder(ctx context.Context, userUUID string, partUUIDs []string) (*model.Order, error)
	GetOrder(ctx context.Context, orderUUID string) (*model.Order, error)
	PayOrder(ctx context.Context, orderUUID string, paymentMethod model.PaymentMethod) (string, error)
	CancelOrder(ctx context.Context, orderUUID string) error
}

type ConsumerService interface {
	RunConsumer(ctx context.Context) error
}

type ProducerService interface {
	ProduceOrderPaid(ctx context.Context, event model.OrderPaidEvent) error
}

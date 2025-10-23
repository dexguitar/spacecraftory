package repository

import (
	"context"

	"github.com/dexguitar/spacecraftory/order/internal/model"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *model.Order) (*model.Order, error)
	GetOrder(ctx context.Context, orderUUID string) (*model.Order, error)
	UpdateOrder(ctx context.Context, order *model.Order) error
}

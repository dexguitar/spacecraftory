package client

import (
	"context"

	"github.com/dexguitar/spacecraftory/order/internal/model"
)

type InventoryClient interface {
	ListParts(ctx context.Context, filter *model.PartsFilter) ([]model.Part, error)
}

type PaymentClient interface {
	PayOrder(ctx context.Context, orderUUID, userUUID string, paymentMethod model.PaymentMethod) (string, error)
}

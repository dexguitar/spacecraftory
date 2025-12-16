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

type IAMClient interface {
	WhoAmI(ctx context.Context, sessionUUID string) (*model.Session, *model.User, error)
}

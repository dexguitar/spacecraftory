package payment

import (
	"context"

	"github.com/dexguitar/spacecraftory/order/internal/model"
)

type PaymentClient interface {
	PayOrder(ctx context.Context, orderUUID, userUUID string, paymentMethod model.PaymentMethod) (string, error)
}

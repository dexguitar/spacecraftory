package repository

import (
	"context"

	"github.com/dexguitar/spacecraftory/payment/internal/model"
)

type PaymentRepository interface {
	PayOrder(ctx context.Context, payment *model.Payment) (string, error)
}

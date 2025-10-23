package service

import (
	"context"

	"github.com/dexguitar/spacecraftory/payment/internal/model"
)

type PaymentService interface {
	PayOrder(ctx context.Context, payment *model.Payment) (string, error)
}

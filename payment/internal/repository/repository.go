package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/dexguitar/spacecraftory/payment/internal/model"
)

type PaymentRepository interface {
	PayOrder(ctx context.Context, payment *model.Payment) (uuid.UUID, error)
}

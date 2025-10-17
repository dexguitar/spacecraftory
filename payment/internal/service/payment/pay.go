package payment

import (
	"context"

	"github.com/dexguitar/spacecraftory/payment/internal/model"
)

func (s *service) PayOrder(ctx context.Context, payment *model.Payment) (string, error) {
	transactionUUID, err := s.paymentRepository.PayOrder(ctx, payment)
	if err != nil {
		// TODO: later will add db error check and map to service errors
		return "", err
	}
	return transactionUUID.String(), nil
}

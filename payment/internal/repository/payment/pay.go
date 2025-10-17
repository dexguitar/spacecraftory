package payment

import (
	"context"

	"github.com/google/uuid"

	"github.com/dexguitar/spacecraftory/payment/internal/model"
	repoConverter "github.com/dexguitar/spacecraftory/payment/internal/repository/converter"
)

func (r *paymentRepository) PayOrder(ctx context.Context, paymentInfo *model.Payment) (uuid.UUID, error) {
	newPaymentUUID := uuid.New()

	repoModel := repoConverter.PaymentInfoToRepoModel(paymentInfo)

	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[newPaymentUUID] = &repoModel

	return newPaymentUUID, nil
}

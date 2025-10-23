package payment

import (
	"github.com/dexguitar/spacecraftory/payment/internal/repository"
)

type service struct {
	paymentRepository repository.PaymentRepository
}

func NewService(paymentRepository repository.PaymentRepository) *service {
	return &service{paymentRepository: paymentRepository}
}

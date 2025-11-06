package app

import (
	"context"

	paymentV1API "github.com/dexguitar/spacecraftory/payment/internal/api/payment/v1"
	"github.com/dexguitar/spacecraftory/payment/internal/repository"
	paymentRepository "github.com/dexguitar/spacecraftory/payment/internal/repository/payment"
	"github.com/dexguitar/spacecraftory/payment/internal/service"
	paymentService "github.com/dexguitar/spacecraftory/payment/internal/service/payment"
	paymentV1 "github.com/dexguitar/spacecraftory/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	paymentV1API paymentV1.PaymentServiceServer

	paymentService service.PaymentService

	paymentRepository repository.PaymentRepository
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) PaymentV1API(ctx context.Context) paymentV1.PaymentServiceServer {
	if d.paymentV1API == nil {
		d.paymentV1API = paymentV1API.NewAPI(d.PaymentService(ctx))
	}

	return d.paymentV1API
}

func (d *diContainer) PaymentService(ctx context.Context) service.PaymentService {
	if d.paymentService == nil {
		d.paymentService = paymentService.NewService(d.PaymentRepository(ctx))
	}

	return d.paymentService
}

func (d *diContainer) PaymentRepository(_ context.Context) repository.PaymentRepository {
	if d.paymentRepository == nil {
		d.paymentRepository = paymentRepository.NewPaymentRepository()
	}

	return d.paymentRepository
}

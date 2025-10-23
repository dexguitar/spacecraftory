package payment

import (
	paymentV1 "github.com/dexguitar/spacecraftory/shared/pkg/proto/payment/v1"
)

type paymentClient struct {
	grpcClient paymentV1.PaymentServiceClient
}

func NewPaymentClient(grpcClient paymentV1.PaymentServiceClient) *paymentClient {
	return &paymentClient{
		grpcClient: grpcClient,
	}
}

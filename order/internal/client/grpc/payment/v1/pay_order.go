package payment

import (
	"context"

	"github.com/dexguitar/spacecraftory/order/internal/client/converter"
	"github.com/dexguitar/spacecraftory/order/internal/model"
	paymentV1 "github.com/dexguitar/spacecraftory/shared/pkg/proto/payment/v1"
)

func (c *paymentClient) PayOrder(ctx context.Context, orderUUID, userUUID string, paymentMethod model.PaymentMethod) (string, error) {
	req := &paymentV1.PayOrderRequest{
		OrderUuid:     orderUUID,
		UserUuid:      userUUID,
		PaymentMethod: converter.PaymentMethodToProto(paymentMethod),
	}

	resp, err := c.grpcClient.PayOrder(ctx, req)
	if err != nil {
		return "", err
	}

	return resp.TransactionUuid, nil
}

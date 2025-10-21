package v1

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/dexguitar/spacecraftory/order/internal/converter"
	"github.com/dexguitar/spacecraftory/order/internal/model"
	orderV1 "github.com/dexguitar/spacecraftory/shared/pkg/openapi/order/v1"
)

func (a *api) PayOrder(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	_, err := uuid.Parse(params.OrderUUID.String())
	if err != nil {
		return &orderV1.BadRequestError{
			Code:    400,
			Message: "Invalid order UUID",
		}, nil
	}
	paymentMethod := converter.ToModelPaymentMethod(req.PaymentMethod)

	transactionUUID, err := a.orderService.PayOrder(ctx, params.OrderUUID.String(), paymentMethod)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return &orderV1.NotFoundError{
				Code:    404,
				Message: "Order not found",
			}, nil
		}
		if errors.Is(err, model.ErrInvalidOrderStatus) {
			return &orderV1.ConflictError{
				Code:    409,
				Message: "Order has already been paid or cancelled",
			}, nil
		}
		return &orderV1.InternalServerError{
			Code:    500,
			Message: "Failed to process payment",
		}, nil
	}

	return &orderV1.PayOrderResponse{
		TransactionUUID: converter.ToProtoTransactionUUID(transactionUUID),
	}, nil
}

package v1

import (
	"context"
	"errors"

	"github.com/dexguitar/spacecraftory/order/internal/model"
	orderV1 "github.com/dexguitar/spacecraftory/shared/pkg/openapi/order/v1"
)

func (a *api) CancelOrder(ctx context.Context, params orderV1.CancelOrderParams) (orderV1.CancelOrderRes, error) {
	err := a.orderService.CancelOrder(ctx, params.OrderUUID.String())
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
				Message: "Cannot cancel already paid or cancelled order",
			}, nil
		}
		return &orderV1.InternalServerError{
			Code:    500,
			Message: "Failed to cancel order",
		}, nil
	}

	return &orderV1.CancelOrderNoContent{}, nil
}

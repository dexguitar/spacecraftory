package v1

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/dexguitar/spacecraftory/order/internal/converter"
	"github.com/dexguitar/spacecraftory/order/internal/model"
	orderV1 "github.com/dexguitar/spacecraftory/shared/pkg/openapi/order/v1"
)

func (a *api) GetOrderByUUID(ctx context.Context, params orderV1.GetOrderByUUIDParams) (orderV1.GetOrderByUUIDRes, error) {
	_, err := uuid.Parse(params.OrderUUID.String())
	if err != nil {
		return &orderV1.BadRequestError{
			Code:    400,
			Message: "Invalid order UUID",
		}, nil
	}
	order, err := a.orderService.GetOrder(ctx, params.OrderUUID.String())
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return &orderV1.NotFoundError{
				Code:    404,
				Message: "Order not found",
			}, nil
		}
		return &orderV1.InternalServerError{
			Code:    500,
			Message: "Failed to get order",
		}, nil
	}

	return converter.ToDtoOrder(order), nil
}

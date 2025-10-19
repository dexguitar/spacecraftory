package v1

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/dexguitar/spacecraftory/order/internal/model"
	orderV1 "github.com/dexguitar/spacecraftory/shared/pkg/openapi/order/v1"
)

func (a *api) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	partUUIDs := make([]string, 0, len(req.PartUuids))
	for _, partUUID := range req.PartUuids {
		partUUIDs = append(partUUIDs, partUUID.String())
	}

	order, err := a.orderService.CreateOrder(ctx, req.UserUUID.String(), partUUIDs)
	if err != nil {
		if errors.Is(err, model.ErrBadRequest) || errors.Is(err, model.ErrPartsNotFound) {
			return &orderV1.BadRequestError{
				Code:    400,
				Message: err.Error(),
			}, nil
		}
		return &orderV1.InternalServerError{
			Code:    500,
			Message: "Failed to create order",
		}, nil
	}

	orderUUID, err := uuid.Parse(order.OrderUUID)
	if err != nil {
		return &orderV1.InternalServerError{
			Code:    500,
			Message: "Failed to parse order UUID",
		}, nil
	}

	return &orderV1.CreateOrderResponse{
		OrderUUID:  orderUUID,
		TotalPrice: order.TotalPrice,
	}, nil
}

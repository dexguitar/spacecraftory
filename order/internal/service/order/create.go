package order

import (
	"context"

	"github.com/dexguitar/spacecraftory/order/internal/model"
)

func (s *service) CreateOrder(ctx context.Context, userUUID string, partUUIDs []string) (*model.Order, error) {
	if len(partUUIDs) == 0 {
		return nil, model.ErrBadRequest
	}

	filter := &model.PartsFilter{
		UUIDs: partUUIDs,
	}
	parts, err := s.inventoryClient.ListParts(ctx, filter)
	if err != nil {
		return nil, err
	}

	if len(parts) == 0 || len(parts) != len(partUUIDs) {
		return nil, model.ErrPartsNotFound
	}

	var totalPrice float64
	for _, part := range parts {
		totalPrice += part.Price
	}

	order := &model.Order{
		UserUUID:   userUUID,
		PartUUIDs:  partUUIDs,
		TotalPrice: totalPrice,
		Status:     model.StatusPENDINGPAYMENT,
	}

	createdOrder, err := s.orderRepository.CreateOrder(ctx, order)
	if err != nil {
		return nil, err
	}

	return createdOrder, nil
}

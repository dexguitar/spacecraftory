package order

import (
	"context"

	"github.com/dexguitar/spacecraftory/order/internal/model"
)

func (s *service) GetOrder(ctx context.Context, orderUUID string) (*model.Order, error) {
	order, err := s.orderRepository.GetOrder(ctx, orderUUID)
	if err != nil {
		return nil, err
	}

	return order, nil
}

package order

import (
	"context"

	"github.com/dexguitar/spacecraftory/order/internal/model"
)

func (s *service) CancelOrder(ctx context.Context, orderUUID string) error {
	order, err := s.orderRepository.GetOrder(ctx, orderUUID)
	if err != nil {
		return err
	}

	if order.OrderStatus == model.OrderStatusPAID || order.OrderStatus == model.OrderStatusCANCELLED {
		return model.ErrInvalidOrderStatus
	}

	order.OrderStatus = model.OrderStatusCANCELLED

	if err := s.orderRepository.UpdateOrder(ctx, order); err != nil {
		return err
	}

	return nil
}

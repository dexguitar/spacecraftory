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

	if order.Status == model.StatusPAID || order.Status == model.StatusCANCELLED {
		return model.ErrInvalidOrderStatus
	}

	order.Status = model.StatusCANCELLED

	if err := s.orderRepository.UpdateOrder(ctx, order); err != nil {
		return err
	}

	return nil
}

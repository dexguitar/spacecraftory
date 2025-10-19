package order

import (
	"context"

	"github.com/dexguitar/spacecraftory/order/internal/model"
)

func (s *service) PayOrder(ctx context.Context, orderUUID string, paymentMethod model.PaymentMethod) (string, error) {
	order, err := s.orderRepository.GetOrder(ctx, orderUUID)
	if err != nil {
		return "", err
	}

	if order.Status == model.StatusPAID || order.Status == model.StatusCANCELLED {
		return "", model.ErrInvalidOrderStatus
	}

	transactionUUID, err := s.paymentClient.PayOrder(ctx, orderUUID, order.UserUUID, paymentMethod)
	if err != nil {
		return "", model.ErrPaymentFailed
	}

	order.Status = model.StatusPAID
	order.TransactionUUID = transactionUUID
	order.PaymentMethod = paymentMethod

	if err := s.orderRepository.UpdateOrder(ctx, order); err != nil {
		return "", err
	}

	return transactionUUID, nil
}

package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/dexguitar/spacecraftory/order/internal/model"
)

func (s *service) PayOrder(ctx context.Context, orderUUID string, paymentMethod model.PaymentMethod) (string, error) {
	order, err := s.orderRepository.GetOrder(ctx, orderUUID)
	if err != nil {
		return "", err
	}

	if order.OrderStatus == model.OrderStatusPAID || order.OrderStatus == model.OrderStatusCANCELLED {
		return "", model.ErrInvalidOrderStatus
	}

	transactionUUID, err := s.paymentClient.PayOrder(ctx, orderUUID, order.UserUUID, paymentMethod)
	if err != nil {
		return "", model.ErrPaymentFailed
	}

	order.OrderStatus = model.OrderStatusPAID
	order.TransactionUUID = transactionUUID
	order.PaymentMethod = paymentMethod

	if err := s.orderRepository.UpdateOrder(ctx, order); err != nil {
		return "", err
	}

	err = s.producerService.ProduceOrderPaid(ctx, model.OrderPaidEvent{
		EventUUID:       uuid.NewString(),
		OrderUUID:       orderUUID,
		UserUUID:        order.UserUUID,
		PaymentMethod:   string(paymentMethod),
		TransactionUUID: transactionUUID,
	})
	if err != nil {
		return "", err
	}

	return transactionUUID, nil
}

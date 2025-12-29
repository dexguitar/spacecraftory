package order

import (
	"context"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/dexguitar/spacecraftory/order/internal/model"
	"github.com/dexguitar/spacecraftory/platform/pkg/tracing"
)

func (s *service) PayOrder(ctx context.Context, orderUUID string, paymentMethod model.PaymentMethod) (string, error) {
	// Create root span for the payment operation
	ctx, span := tracing.StartSpan(ctx, "order.PayOrder",
		trace.WithAttributes(
			attribute.String("order.uuid", orderUUID),
			attribute.String("order.payment_method", string(paymentMethod)),
		),
	)
	defer span.End()

	order, err := s.orderRepository.GetOrder(ctx, orderUUID)
	if err != nil {
		span.RecordError(err)
		return "", err
	}

	if order.OrderStatus == model.OrderStatusPAID || order.OrderStatus == model.OrderStatusCANCELLED {
		span.RecordError(model.ErrInvalidOrderStatus)
		return "", model.ErrInvalidOrderStatus
	}

	// Add order details to span
	span.SetAttributes(
		attribute.String("order.user_uuid", order.UserUUID),
		attribute.Float64("order.total_price", order.TotalPrice),
	)

	transactionUUID, err := s.paymentClient.PayOrder(ctx, orderUUID, order.UserUUID, paymentMethod)
	if err != nil {
		span.RecordError(err)
		return "", model.ErrPaymentFailed
	}

	order.OrderStatus = model.OrderStatusPAID
	order.TransactionUUID = transactionUUID
	order.PaymentMethod = paymentMethod

	if err := s.orderRepository.UpdateOrder(ctx, order); err != nil {
		span.RecordError(err)
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
		span.RecordError(err)
		return "", err
	}

	// Add success attributes
	span.SetAttributes(
		attribute.String("order.transaction_uuid", transactionUUID),
		attribute.String("order.status", string(model.OrderStatusPAID)),
	)

	return transactionUUID, nil
}

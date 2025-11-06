package order

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	serviceModel "github.com/dexguitar/spacecraftory/order/internal/model"
)

func (r *orderRepository) UpdateOrder(ctx context.Context, order *serviceModel.Order) error {
	builderUpdate := sq.
		Update("orders").
		PlaceholderFormat(sq.Dollar).
		Set("status", order.OrderStatus).
		Set("transaction_uuid", order.TransactionUUID).
		Set("payment_method", order.PaymentMethod).
		Where(sq.Eq{"id": order.OrderUUID})

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

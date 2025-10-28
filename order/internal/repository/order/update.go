package order

import (
	"context"
	"log"

	sq "github.com/Masterminds/squirrel"

	serviceModel "github.com/dexguitar/spacecraftory/order/internal/model"
)

func (r *orderRepository) UpdateOrder(ctx context.Context, order *serviceModel.Order) error {
	builderUpdate := sq.Update("orders").PlaceholderFormat(sq.Dollar).Set("status", order.OrderStatus).Set("transaction_uuid", order.TransactionUUID).Set("payment_method", order.PaymentMethod).Where(sq.Eq{"id": order.OrderUUID})
	query, args, err := builderUpdate.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v\n", serviceModel.ErrInternalServerError)
		return err
	}

	res, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		log.Printf("failed to update order: %v\n", err)
		return err
	}

	log.Printf("deleted %d rows", res.RowsAffected())

	return nil
}

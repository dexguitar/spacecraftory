package order

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	serviceModel "github.com/dexguitar/spacecraftory/order/internal/model"
	"github.com/dexguitar/spacecraftory/order/internal/repository/converter"
	"github.com/dexguitar/spacecraftory/order/internal/repository/model"
)

func (r *orderRepository) GetOrder(ctx context.Context, orderUUID string) (*serviceModel.Order, error) {
	builderSelectOne := sq.
		Select("id", "user_uuid", "part_uuids", "total_price", "status", "transaction_uuid", "payment_method").
		From("orders").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": orderUUID}).
		Limit(1)

	query, args, err := builderSelectOne.ToSql()
	if err != nil {
		return nil, err
	}

	var order model.Order
	err = r.db.QueryRow(ctx, query, args...).Scan(&order.OrderUUID, &order.UserUUID, &order.PartUUIDs, &order.TotalPrice, &order.Status, &order.TransactionUUID, &order.PaymentMethod)
	if err != nil {
		return nil, err
	}

	return converter.ToModelOrder(&order), nil
}

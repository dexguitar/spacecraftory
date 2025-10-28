package order

import (
	"context"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	serviceModel "github.com/dexguitar/spacecraftory/order/internal/model"
	repoConverter "github.com/dexguitar/spacecraftory/order/internal/repository/converter"
)

func (r *orderRepository) CreateOrder(ctx context.Context, order *serviceModel.Order) (*serviceModel.Order, error) {
	repoOrder := repoConverter.ToRepoOrder(order)

	builderInsert := sq.Insert("orders").PlaceholderFormat(sq.Dollar).Columns("user_uuid", "part_uuids", "total_price", "payment_method").Values(repoOrder.UserUUID, repoOrder.PartUUIDs, repoOrder.TotalPrice, repoOrder.PaymentMethod).Suffix("RETURNING id")
	query, args, err := builderInsert.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v\n", serviceModel.ErrInternalServerError)
		return nil, err
	}

	var orderID uuid.UUID
	err = r.db.QueryRow(ctx, query, args...).Scan(&orderID)
	if err != nil {
		log.Printf("failed to create order: %v\n", err)
		return nil, err
	}

	log.Printf("created order with id: %d", orderID)

	repoOrder.OrderUUID = orderID.String()

	return repoConverter.ToModelOrder(repoOrder), nil
}

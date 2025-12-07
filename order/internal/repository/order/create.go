package order

import (
	"context"
	"errors"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"

	serviceModel "github.com/dexguitar/spacecraftory/order/internal/model"
)

type orderID struct {
	ID string `db:"id"`
}

func (r *orderRepository) CreateOrder(ctx context.Context, order *serviceModel.Order) (*serviceModel.Order, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = tx.Rollback(ctx)
		if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			log.Printf("failed to rollback transaction: %v", err)
		}
	}()

	// Create order in orders table
	orderInsert := sq.Insert("orders").
		PlaceholderFormat(sq.Dollar).
		Columns("user_uuid", "total_price", "status").
		Values(order.UserUUID, order.TotalPrice, order.OrderStatus).
		Suffix("RETURNING id")

	query, args, err := orderInsert.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[orderID])
	if err != nil {
		return nil, err
	}

	// Insert parts into order_parts table
	if len(order.PartUUIDs) > 0 {
		partsInsert := sq.Insert("order_parts").
			PlaceholderFormat(sq.Dollar).
			Columns("order_id", "part_id")

		for _, partUUID := range order.PartUUIDs {
			partsInsert = partsInsert.Values(result.ID, partUUID)
		}

		partsQuery, partsArgs, err := partsInsert.ToSql()
		if err != nil {
			return nil, err
		}

		_, err = tx.Exec(ctx, partsQuery, partsArgs...)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	order.OrderUUID = result.ID
	return order, nil
}

package order

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"

	serviceModel "github.com/dexguitar/spacecraftory/order/internal/model"
	"github.com/dexguitar/spacecraftory/order/internal/repository/converter"
	"github.com/dexguitar/spacecraftory/order/internal/repository/model"
)

type orderPart struct {
	PartID string `db:"part_id"`
}

func (r *orderRepository) GetOrder(ctx context.Context, orderUUID string) (*serviceModel.Order, error) {
	// Get order
	orderQuery := sq.
		Select("id", "user_uuid", "total_price", "status", "transaction_uuid", "payment_method").
		From("orders").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": orderUUID})

	query, args, err := orderQuery.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	order, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Order])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, serviceModel.ErrOrderNotFound
		}
		return nil, err
	}

	// Get order parts
	partsQuery := sq.Select("part_id").
		From("order_parts").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"order_id": orderUUID})

	partsQueryStr, partsArgs, err := partsQuery.ToSql()
	if err != nil {
		return nil, err
	}

	partRows, err := r.db.Query(ctx, partsQueryStr, partsArgs...)
	if err != nil {
		return nil, err
	}
	defer partRows.Close()

	parts, err := pgx.CollectRows(partRows, pgx.RowToStructByName[orderPart])
	if err != nil {
		return nil, err
	}

	partUUIDs := make([]string, len(parts))
	for i, part := range parts {
		partUUIDs[i] = part.PartID
	}

	serviceOrder := converter.ToModelOrder(&order)
	serviceOrder.PartUUIDs = partUUIDs

	return serviceOrder, nil
}

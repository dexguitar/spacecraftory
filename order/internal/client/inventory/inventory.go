package inventory

import (
	"context"

	"github.com/dexguitar/spacecraftory/order/internal/model"
)

type InventoryClient interface {
	ListParts(ctx context.Context, filter *model.PartsFilter) ([]model.Part, error)
}

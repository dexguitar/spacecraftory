package converter

import (
	"github.com/dexguitar/spacecraftory/order/internal/model"
	inventoryV1 "github.com/dexguitar/spacecraftory/shared/pkg/proto/inventory/v1"
)

func PartsFilterToProto(filter *model.PartsFilter) *inventoryV1.PartsFilter {
	if filter == nil {
		return nil
	}

	return &inventoryV1.PartsFilter{
		Uuids: filter.UUIDs,
	}
}

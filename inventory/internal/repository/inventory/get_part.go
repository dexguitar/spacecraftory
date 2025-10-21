package inventory

import (
	"context"

	"github.com/dexguitar/spacecraftory/inventory/internal/model"
	repoConverter "github.com/dexguitar/spacecraftory/inventory/internal/repository/converter"
)

func (r *inventoryRepository) GetPart(ctx context.Context, uuid string) (*model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	repoPart, ok := r.parts[uuid]
	if !ok {
		return nil, model.ErrPartNotFound
	}

	return repoConverter.ToModelPart(repoPart), nil
}

package inventory

import (
	"context"

	"github.com/dexguitar/spacecraftory/inventory/internal/model"
)

func (s *service) ListParts(ctx context.Context, filter *model.PartsFilter) ([]*model.Part, error) {
	parts, err := s.inventoryRepository.ListParts(ctx, filter)
	if err != nil {
		// TODO: add business logic validation or error mapping
		return nil, err
	}

	return parts, nil
}

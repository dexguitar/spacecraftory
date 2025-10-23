package inventory

import (
	"context"

	"github.com/dexguitar/spacecraftory/inventory/internal/model"
)

func (s *service) GetPart(ctx context.Context, uuid string) (*model.Part, error) {
	part, err := s.inventoryRepository.GetPart(ctx, uuid)
	if err != nil {
		return nil, err
	}

	return part, nil
}

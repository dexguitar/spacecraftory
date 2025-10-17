package inventory

import (
	"sync"

	repoModel "github.com/dexguitar/spacecraftory/inventory/internal/repository/model"
)

type inventoryRepository struct {
	mu    sync.RWMutex
	parts map[string]*repoModel.Part
}

func NewInventoryRepository() *inventoryRepository {
	return &inventoryRepository{
		parts: make(map[string]*repoModel.Part),
	}
}

func (r *inventoryRepository) InitializeMockData(parts map[string]*repoModel.Part) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.parts = parts
}

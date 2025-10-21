package inventory

import (
	"sync"
	"time"

	"github.com/dexguitar/spacecraftory/inventory/internal/model"
	repoModel "github.com/dexguitar/spacecraftory/inventory/internal/repository/model"
)

type inventoryRepository struct {
	mu    sync.RWMutex
	parts map[string]*repoModel.Part
}

func NewInventoryRepository() *inventoryRepository {
	return &inventoryRepository{
		parts: initParts(),
	}
}

func initParts() map[string]*repoModel.Part {
	now := time.Now()

	mockParts := map[string]*repoModel.Part{
		"123e4567-e89b-12d3-a456-426614174000": {
			UUID:          "123e4567-e89b-12d3-a456-426614174000",
			Name:          "Quantum Drive Engine",
			Description:   "High-efficiency quantum propulsion engine for interstellar travel",
			Price:         150000.00,
			StockQuantity: 5,
			Category:      model.CategoryEngine,
			Dimensions: &repoModel.Dimensions{
				Length: 3.5,
				Width:  2.0,
				Height: 2.5,
				Weight: 500.0,
			},
			Manufacturer: &repoModel.Manufacturer{
				Name:    "SpaceTech Industries",
				Country: "USA",
				Website: "https://spacetech.example.com",
			},
			Tags:      []string{"quantum", "propulsion", "interstellar"},
			CreatedAt: now,
			UpdatedAt: now,
		},
		"123e4567-e89b-12d3-a456-426614174001": {
			UUID:          "123e4567-e89b-12d3-a456-426614174001",
			Name:          "Fusion Fuel Cell",
			Description:   "Advanced fusion-based fuel cell for long-duration missions",
			Price:         75000.00,
			StockQuantity: 12,
			Category:      model.CategoryFuel,
			Dimensions: &repoModel.Dimensions{
				Length: 1.2,
				Width:  0.8,
				Height: 1.0,
				Weight: 100.0,
			},
			Manufacturer: &repoModel.Manufacturer{
				Name:    "Energy Solutions Corp",
				Country: "Germany",
				Website: "https://energysolutions.example.com",
			},
			Tags:      []string{"fusion", "fuel", "efficient"},
			CreatedAt: now,
			UpdatedAt: now,
		},
		"123e4567-e89b-12d3-a456-426614174002": {
			UUID:          "123e4567-e89b-12d3-a456-426614174002",
			Name:          "Reinforced Porthole",
			Description:   "Triple-layered reinforced viewing porthole for crew observation",
			Price:         25000.00,
			StockQuantity: 20,
			Category:      model.CategoryPorthole,
			Dimensions: &repoModel.Dimensions{
				Length: 1.0,
				Width:  1.0,
				Height: 0.3,
				Weight: 50.0,
			},
			Manufacturer: &repoModel.Manufacturer{
				Name:    "ViewTech Manufacturing",
				Country: "Japan",
				Website: "https://viewtech.example.com",
			},
			Tags:      []string{"observation", "reinforced", "safety"},
			CreatedAt: now,
			UpdatedAt: now,
		},
		"123e4567-e89b-12d3-a456-426614174003": {
			UUID:          "123e4567-e89b-12d3-a456-426614174003",
			Name:          "Aerodynamic Wing Panel",
			Description:   "Carbon-fiber composite wing panel for atmospheric re-entry",
			Price:         45000.00,
			StockQuantity: 8,
			Category:      model.CategoryWing,
			Dimensions: &repoModel.Dimensions{
				Length: 5.0,
				Width:  2.5,
				Height: 0.5,
				Weight: 200.0,
			},
			Manufacturer: &repoModel.Manufacturer{
				Name:    "AeroDynamics Ltd",
				Country: "UK",
				Website: "https://aerodynamics.example.com",
			},
			Tags:      []string{"aerodynamic", "reentry", "composite"},
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	return mockParts
}

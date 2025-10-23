package v1

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/dexguitar/spacecraftory/inventory/internal/model"
	"github.com/dexguitar/spacecraftory/inventory/internal/service/mocks"
)

type APISuite struct {
	suite.Suite
	ctx              context.Context
	inventoryService *mocks.InventoryService
	api              *api
	serviceMockData  map[string]*model.Part
}

func (s *APISuite) SetupTest() {
	s.ctx = context.Background()
	s.inventoryService = mocks.NewInventoryService(s.T())
	s.api = NewAPI(s.inventoryService)
	s.serviceMockData = generateServiceMockData()
}

func TestAPIIntegration(t *testing.T) {
	suite.Run(t, new(APISuite))
}

func generateServiceMockData() map[string]*model.Part {
	now := time.Now()

	return map[string]*model.Part{
		"123e4567-e89b-12d3-a456-426614174000": {
			UUID:          "123e4567-e89b-12d3-a456-426614174000",
			Name:          "Quantum Drive Engine",
			Description:   "High-efficiency quantum propulsion engine",
			Price:         150000.00,
			StockQuantity: 5,
			Category:      model.CategoryEngine,
			Dimensions: &model.Dimensions{
				Length: 3.5,
				Width:  2.0,
				Height: 2.5,
				Weight: 500.0,
			},
			Manufacturer: &model.Manufacturer{
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
			Description:   "Advanced fusion-based fuel cell",
			Price:         75000.00,
			StockQuantity: 12,
			Category:      model.CategoryFuel,
			Dimensions: &model.Dimensions{
				Length: 1.2,
				Width:  0.8,
				Height: 1.0,
				Weight: 100.0,
			},
			Manufacturer: &model.Manufacturer{
				Name:    "Energy Solutions Corp",
				Country: "Germany",
				Website: "https://energysolutions.example.com",
			},
			Tags:      []string{"fusion", "fuel", "long-duration"},
			CreatedAt: now,
			UpdatedAt: now,
		},
		"123e4567-e89b-12d3-a456-426614174002": {
			UUID:          "123e4567-e89b-12d3-a456-426614174002",
			Name:          "Reinforced Porthole",
			Description:   "Triple-layered reinforced viewing porthole",
			Price:         25000.00,
			StockQuantity: 20,
			Category:      model.CategoryPorthole,
			Dimensions: &model.Dimensions{
				Length: 0.8,
				Width:  0.8,
				Height: 0.3,
				Weight: 50.0,
			},
			Manufacturer: &model.Manufacturer{
				Name:    "ViewTech Manufacturing",
				Country: "Japan",
				Website: "https://viewtech.example.com",
			},
			Tags:      []string{"porthole", "reinforced", "observation"},
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
}

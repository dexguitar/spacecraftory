package inventory

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/dexguitar/spacecraftory/inventory/internal/model"
	"github.com/dexguitar/spacecraftory/inventory/internal/repository/mocks"
	repoModel "github.com/dexguitar/spacecraftory/inventory/internal/repository/model"
)

type ServiceSuite struct {
	suite.Suite

	ctx context.Context

	inventoryRepo *mocks.InventoryRepository

	service *service

	repoMockData map[string]*repoModel.Part
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()

	s.inventoryRepo = mocks.NewInventoryRepository(s.T())

	s.service = NewService(
		s.inventoryRepo,
	)

	s.repoMockData = generateRepoMockData()
}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}

func generateRepoMockData() map[string]*repoModel.Part {
	return map[string]*repoModel.Part{
		"123e4567-e89b-12d3-a456-426614174000": {
			UUID:          "123e4567-e89b-12d3-a456-426614174000",
			Name:          "Quantum Drive Engine",
			Description:   "High-efficiency quantum propulsion engine for interstellar travel",
			Price:         150000.00,
			StockQuantity: 5,
			Category:      model.CategoryEngine,
		},
		"123e4567-e89b-12d3-a456-426614174001": {
			UUID:          "123e4567-e89b-12d3-a456-426614174001",
			Name:          "Fusion Fuel Cell",
			Description:   "Advanced fusion-based fuel cell for long-duration missions",
			Price:         75000.00,
			StockQuantity: 12,
			Category:      model.CategoryFuel,
		},
		"123e4567-e89b-12d3-a456-426614174002": {
			UUID:          "123e4567-e89b-12d3-a456-426614174002",
			Name:          "Aerodynamic Wing Panel",
			Description:   "Carbon-fiber composite wing panel for atmospheric re-entry",
			Price:         45000.00,
			StockQuantity: 8,
			Category:      model.CategoryWing,
		},
		"123e4567-e89b-12d3-a456-426614174003": {
			UUID:          "123e4567-e89b-12d3-a456-426614174003",
			Name:          "Porthole Window",
			Description:   "Glass porthole window for viewing the stars",
			Price:         10000.00,
			StockQuantity: 10,
			Category:      model.CategoryPorthole,
		},
	}
}

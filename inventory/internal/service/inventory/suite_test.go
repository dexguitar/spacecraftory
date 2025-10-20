package inventory

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/dexguitar/spacecraftory/inventory/internal/repository/mocks"
	repoModel "github.com/dexguitar/spacecraftory/inventory/internal/repository/model"
	"github.com/dexguitar/spacecraftory/inventory/internal/repository/utils"
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

	s.repoMockData = utils.GenerateMockParts()
}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}

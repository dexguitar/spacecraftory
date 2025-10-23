package inventory

import (
	"github.com/stretchr/testify/assert"

	"github.com/dexguitar/spacecraftory/inventory/internal/model"
	"github.com/dexguitar/spacecraftory/inventory/internal/repository/converter"
)

func (s *ServiceSuite) TestGetPartSuccess() {
	// using UUIDs from generated mock data
	validUUIDs := make([]string, 0, len(s.repoMockData))
	for uuid := range s.repoMockData {
		validUUIDs = append(validUUIDs, uuid)
	}

	for _, uuid := range validUUIDs {
		repoPart := s.repoMockData[uuid]
		expectedServicePart := converter.ToModelPart(repoPart)

		s.inventoryRepo.On("GetPart", s.ctx, uuid).Return(expectedServicePart, nil).Once()

		part, err := s.service.GetPart(s.ctx, uuid)

		s.Require().NoError(err)
		assert.Equal(s.T(), expectedServicePart.UUID, part.UUID)
		assert.Equal(s.T(), expectedServicePart.Name, part.Name)
		assert.Equal(s.T(), expectedServicePart.Description, part.Description)
		assert.Equal(s.T(), expectedServicePart.Price, part.Price)
		assert.Equal(s.T(), expectedServicePart.StockQuantity, part.StockQuantity)
		assert.Equal(s.T(), expectedServicePart.Category, part.Category)
		assert.Equal(s.T(), expectedServicePart.Dimensions, part.Dimensions)
		assert.Equal(s.T(), expectedServicePart.Manufacturer, part.Manufacturer)
		assert.Equal(s.T(), expectedServicePart.Tags, part.Tags)
	}
}

func (s *ServiceSuite) TestGetPartError() {
	nonExistentUUID := "non-existent-uuid"

	s.inventoryRepo.On("GetPart", s.ctx, nonExistentUUID).
		Return(nil, model.ErrPartNotFound).Once()

	part, err := s.service.GetPart(s.ctx, nonExistentUUID)

	assert.ErrorIs(s.T(), err, model.ErrPartNotFound)
	assert.Nil(s.T(), part)
	s.inventoryRepo.AssertExpectations(s.T())
}

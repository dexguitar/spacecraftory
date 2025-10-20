package inventory

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/dexguitar/spacecraftory/inventory/internal/model"
	"github.com/dexguitar/spacecraftory/inventory/internal/repository/converter"
)

func (s *ServiceSuite) TestListPartsNoFilterSuccess() {
	expectedParts := make([]*model.Part, 0, len(s.repoMockData))
	for _, repoPart := range s.repoMockData {
		expectedParts = append(expectedParts, converter.PartRepoToServiceModel(repoPart))
	}

	s.inventoryRepo.On("ListParts", s.ctx, mock.Anything).
		Return(expectedParts, nil).Once()

	parts, err := s.service.ListParts(s.ctx, nil)

	s.Require().NoError(err)
	assert.Len(s.T(), parts, len(s.repoMockData))
}

func (s *ServiceSuite) TestListPartsSuccess() {
	testCases := []struct {
		name          string
		filter        *model.PartsFilter
		repoReturn    func() []*model.Part
		expectedParts []*model.Part
	}{
		{
			name: "Category filter",
			filter: &model.PartsFilter{
				Categories: []model.Category{model.CategoryEngine, model.CategoryFuel},
			},
			repoReturn: func() []*model.Part {
				return []*model.Part{
					converter.PartRepoToServiceModel(s.repoMockData["123e4567-e89b-12d3-a456-426614174000"]),
					converter.PartRepoToServiceModel(s.repoMockData["123e4567-e89b-12d3-a456-426614174001"]),
				}
			},
			expectedParts: []*model.Part{
				converter.PartRepoToServiceModel(s.repoMockData["123e4567-e89b-12d3-a456-426614174000"]),
				converter.PartRepoToServiceModel(s.repoMockData["123e4567-e89b-12d3-a456-426614174001"]),
			},
		},
		{
			name: "Manufacturer country filter",
			filter: &model.PartsFilter{
				ManufacturerCountries: []string{"USA"},
			},
			repoReturn: func() []*model.Part {
				return []*model.Part{
					converter.PartRepoToServiceModel(s.repoMockData["123e4567-e89b-12d3-a456-426614174000"]),
				}
			},
			expectedParts: []*model.Part{
				converter.PartRepoToServiceModel(s.repoMockData["123e4567-e89b-12d3-a456-426614174000"]),
			},
		},
		{
			name: "Tag filter",
			filter: &model.PartsFilter{
				Tags: []string{"quantum", "fusion"},
			},
			repoReturn: func() []*model.Part {
				return []*model.Part{
					converter.PartRepoToServiceModel(s.repoMockData["123e4567-e89b-12d3-a456-426614174000"]),
					converter.PartRepoToServiceModel(s.repoMockData["123e4567-e89b-12d3-a456-426614174001"]),
					converter.PartRepoToServiceModel(s.repoMockData["123e4567-e89b-12d3-a456-426614174002"]),
				}
			},
			expectedParts: []*model.Part{
				converter.PartRepoToServiceModel(s.repoMockData["123e4567-e89b-12d3-a456-426614174000"]),
				converter.PartRepoToServiceModel(s.repoMockData["123e4567-e89b-12d3-a456-426614174001"]),
				converter.PartRepoToServiceModel(s.repoMockData["123e4567-e89b-12d3-a456-426614174002"]),
			},
		},
		{
			name: "Name filter",
			filter: &model.PartsFilter{
				Names: []string{"Part 1", "Part 2"},
			},
			repoReturn: func() []*model.Part {
				return []*model.Part{
					converter.PartRepoToServiceModel(s.repoMockData["123e4567-e89b-12d3-a456-426614174000"]),
					converter.PartRepoToServiceModel(s.repoMockData["123e4567-e89b-12d3-a456-426614174001"]),
				}
			},
			expectedParts: []*model.Part{
				converter.PartRepoToServiceModel(s.repoMockData["123e4567-e89b-12d3-a456-426614174000"]),
				converter.PartRepoToServiceModel(s.repoMockData["123e4567-e89b-12d3-a456-426614174001"]),
			},
		},
		{
			name: "UUID filter",
			filter: &model.PartsFilter{
				UUIDs: []string{"123e4567-e89b-12d3-a456-426614174000", "123e4567-e89b-12d3-a456-426614174005"},
			},
			repoReturn: func() []*model.Part {
				return []*model.Part{
					converter.PartRepoToServiceModel(s.repoMockData["123e4567-e89b-12d3-a456-426614174000"]),
				}
			},
			expectedParts: []*model.Part{
				converter.PartRepoToServiceModel(s.repoMockData["123e4567-e89b-12d3-a456-426614174000"]),
			},
		},
		{
			name: "All filters",
			filter: &model.PartsFilter{
				Categories:            []model.Category{model.CategoryEngine, model.CategoryFuel},
				ManufacturerCountries: []string{"USA", "Canada"},
				Tags:                  []string{"tag1", "tag2"},
				Names:                 []string{"Part 1", "Part 2"},
			},
			repoReturn: func() []*model.Part {
				return []*model.Part{
					converter.PartRepoToServiceModel(s.repoMockData["123e4567-e89b-12d3-a456-426614174000"]),
					converter.PartRepoToServiceModel(s.repoMockData["123e4567-e89b-12d3-a456-426614174001"]),
				}
			},
			expectedParts: []*model.Part{
				converter.PartRepoToServiceModel(s.repoMockData["123e4567-e89b-12d3-a456-426614174000"]),
				converter.PartRepoToServiceModel(s.repoMockData["123e4567-e89b-12d3-a456-426614174001"]),
			},
		},
		{
			name:   "No filters",
			filter: nil,
			repoReturn: func() []*model.Part {
				return []*model.Part{
					converter.PartRepoToServiceModel(s.repoMockData["123e4567-e89b-12d3-a456-426614174000"]),
					converter.PartRepoToServiceModel(s.repoMockData["123e4567-e89b-12d3-a456-426614174001"]),
					converter.PartRepoToServiceModel(s.repoMockData["123e4567-e89b-12d3-a456-426614174002"]),
					converter.PartRepoToServiceModel(s.repoMockData["123e4567-e89b-12d3-a456-426614174003"]),
				}
			},
			expectedParts: []*model.Part{
				converter.PartRepoToServiceModel(s.repoMockData["123e4567-e89b-12d3-a456-426614174000"]),
				converter.PartRepoToServiceModel(s.repoMockData["123e4567-e89b-12d3-a456-426614174001"]),
				converter.PartRepoToServiceModel(s.repoMockData["123e4567-e89b-12d3-a456-426614174002"]),
				converter.PartRepoToServiceModel(s.repoMockData["123e4567-e89b-12d3-a456-426614174003"]),
			},
		},
		{
			name: "Empty result",
			filter: &model.PartsFilter{
				UUIDs: []string{"non-existent-uuid"},
			},
			repoReturn: func() []*model.Part {
				return []*model.Part{}
			},
			expectedParts: []*model.Part{},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.inventoryRepo.On("ListParts", s.ctx, tc.filter).
				Return(tc.repoReturn(), nil).Once()

			parts, err := s.service.ListParts(s.ctx, tc.filter)

			s.Require().NoError(err)
			assert.Len(s.T(), parts, len(tc.expectedParts))
			assert.ElementsMatch(s.T(), parts, tc.expectedParts)
		})
	}
}

func (s *ServiceSuite) TestListPartsError() {
	testCases := []struct {
		name      string
		filter    *model.PartsFilter
		repoError error
	}{
		{
			name: "Repository error with filter",
			filter: &model.PartsFilter{
				Categories: []model.Category{model.CategoryEngine},
			},
			repoError: model.ErrPartNotFound,
		},
		{
			name:      "Repository error without filter",
			filter:    nil,
			repoError: assert.AnError,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.inventoryRepo.On("ListParts", s.ctx, mock.Anything).
				Return(nil, tc.repoError).Once()

			parts, err := s.service.ListParts(s.ctx, tc.filter)

			assert.ErrorIs(s.T(), err, tc.repoError)
			assert.Nil(s.T(), parts)
		})
	}
}

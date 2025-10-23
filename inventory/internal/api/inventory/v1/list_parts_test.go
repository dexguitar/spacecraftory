package v1

import (
	"errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/dexguitar/spacecraftory/inventory/internal/converter"
	"github.com/dexguitar/spacecraftory/inventory/internal/model"
	inventoryV1 "github.com/dexguitar/spacecraftory/shared/pkg/proto/inventory/v1"
)

func (s *APISuite) TestListPartsSuccess() {
	testCases := []struct {
		name          string
		request       *inventoryV1.ListPartsRequest
		serviceReturn []*model.Part
	}{
		{
			name:    "No filter",
			request: &inventoryV1.ListPartsRequest{},
			serviceReturn: []*model.Part{
				s.serviceMockData["123e4567-e89b-12d3-a456-426614174000"],
			},
		},
		{
			name: "Category filter",
			request: &inventoryV1.ListPartsRequest{
				Filter: &inventoryV1.PartsFilter{
					Categories: []inventoryV1.Category{
						inventoryV1.Category_CATEGORY_ENGINE,
						inventoryV1.Category_CATEGORY_FUEL,
					},
				},
			},
			serviceReturn: []*model.Part{
				s.serviceMockData["123e4567-e89b-12d3-a456-426614174000"],
				s.serviceMockData["123e4567-e89b-12d3-a456-426614174001"],
			},
		},
		{
			name: "Manufacturer country filter",
			request: &inventoryV1.ListPartsRequest{
				Filter: &inventoryV1.PartsFilter{
					ManufacturerCountries: []string{"USA"},
				},
			},
			serviceReturn: []*model.Part{
				s.serviceMockData["123e4567-e89b-12d3-a456-426614174000"],
			},
		},
		{
			name: "Tag filter",
			request: &inventoryV1.ListPartsRequest{
				Filter: &inventoryV1.PartsFilter{
					Tags: []string{"quantum", "fusion"},
				},
			},
			serviceReturn: []*model.Part{
				s.serviceMockData["123e4567-e89b-12d3-a456-426614174000"],
				s.serviceMockData["123e4567-e89b-12d3-a456-426614174001"],
			},
		},
		{
			name: "UUID filter",
			request: &inventoryV1.ListPartsRequest{
				Filter: &inventoryV1.PartsFilter{
					Uuids: []string{"123e4567-e89b-12d3-a456-426614174000"},
				},
			},
			serviceReturn: []*model.Part{
				s.serviceMockData["123e4567-e89b-12d3-a456-426614174000"],
			},
		},
		{
			name: "All filters",
			request: &inventoryV1.ListPartsRequest{
				Filter: &inventoryV1.PartsFilter{
					Names:                 []string{"Quantum Drive Engine", "Fusion Fuel Cell"},
					Categories:            []inventoryV1.Category{inventoryV1.Category_CATEGORY_ENGINE},
					ManufacturerCountries: []string{"USA", "Germany"},
					Tags:                  []string{"quantum", "fusion"},
				},
			},
			serviceReturn: []*model.Part{
				s.serviceMockData["123e4567-e89b-12d3-a456-426614174000"],
				s.serviceMockData["123e4567-e89b-12d3-a456-426614174001"],
				s.serviceMockData["123e4567-e89b-12d3-a456-426614174002"],
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.inventoryService.On("ListParts", s.ctx, mock.Anything).
				Return(tc.serviceReturn, nil).Once()

			resp, err := s.api.ListParts(s.ctx, tc.request)

			s.Require().NoError(err)
			assert.Len(s.T(), resp.Parts, len(tc.serviceReturn))

			expectedProtoParts := make([]*inventoryV1.Part, 0, len(tc.serviceReturn))
			for _, servicePart := range tc.serviceReturn {
				expectedProtoParts = append(expectedProtoParts, converter.ToProtoPart(servicePart))
			}

			assert.ElementsMatch(s.T(), expectedProtoParts, resp.Parts)
		})
	}
}

func (s *APISuite) TestListPartsError() {
	testCases := []struct {
		name             string
		request          *inventoryV1.ListPartsRequest
		serviceError     error
		expectedMsgParts []string
	}{
		{
			name: "Internal error with filter",
			request: &inventoryV1.ListPartsRequest{
				Filter: &inventoryV1.PartsFilter{
					Categories: []inventoryV1.Category{inventoryV1.Category_CATEGORY_ENGINE},
				},
			},
			serviceError:     errors.New("some error"),
			expectedMsgParts: []string{"Internal server error"},
		},
		{
			name: "Internal error without filter",
			request: &inventoryV1.ListPartsRequest{
				Filter: nil,
			},
			serviceError:     errors.New("some error"),
			expectedMsgParts: []string{"Internal server error"},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.inventoryService.On("ListParts", s.ctx, mock.Anything).
				Return(nil, tc.serviceError).Once()

			resp, err := s.api.ListParts(s.ctx, tc.request)

			s.Require().Error(err)
			s.Require().Nil(resp)

			st, ok := status.FromError(err)
			s.Require().True(ok)
			assert.Equal(s.T(), codes.Internal, st.Code())
			for _, msgPart := range tc.expectedMsgParts {
				assert.Contains(s.T(), st.Message(), msgPart)
			}
		})
	}
}

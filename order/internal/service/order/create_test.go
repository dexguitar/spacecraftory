package order

import (
	"errors"

	"github.com/stretchr/testify/assert"

	"github.com/dexguitar/spacecraftory/order/internal/model"
)

var (
	ErrInventoryServiceError = errors.New("inventory service error")
	ErrCreateOrderError      = errors.New("create order error")
)

func (s *OrderServiceSuite) TestCreateOrderSuccess() {
	testCases := []struct {
		name          string
		userUUID      string
		partUUIDs     []string
		parts         []model.Part
		expectedPrice float64
	}{
		{
			name:      "Single part",
			userUUID:  "123e4567-e89b-12d3-a456-426614174000",
			partUUIDs: []string{"part-uuid-1"},
			parts: []model.Part{
				{
					UUID:  "part-uuid-1",
					Name:  "Engine",
					Price: 100.50,
				},
			},
			expectedPrice: 100.50,
		},
		{
			name:      "Multiple parts",
			userUUID:  "123e4567-e89b-12d3-a456-426614174000",
			partUUIDs: []string{"part-uuid-1", "part-uuid-2", "part-uuid-3"},
			parts: []model.Part{
				{
					UUID:  "part-uuid-1",
					Name:  "Engine",
					Price: 100.50,
				},
				{
					UUID:  "part-uuid-2",
					Name:  "Fuel Cell",
					Price: 75.25,
				},
				{
					UUID:  "part-uuid-3",
					Name:  "Wing",
					Price: 50.00,
				},
			},
			expectedPrice: 225.75,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			filter := &model.PartsFilter{
				UUIDs: tc.partUUIDs,
			}

			s.inventoryClient.On("ListParts", s.ctx, filter).
				Return(tc.parts, nil).Once()

			expectedOrder := &model.Order{
				UserUUID:    tc.userUUID,
				PartUUIDs:   tc.partUUIDs,
				TotalPrice:  tc.expectedPrice,
				OrderStatus: model.OrderStatusPENDINGPAYMENT,
			}

			createdOrder := &model.Order{
				OrderUUID:   "created-order-uuid",
				UserUUID:    tc.userUUID,
				PartUUIDs:   tc.partUUIDs,
				TotalPrice:  tc.expectedPrice,
				OrderStatus: model.OrderStatusPENDINGPAYMENT,
			}

			s.orderRepository.On("CreateOrder", s.ctx, expectedOrder).
				Return(createdOrder, nil).Once()

			order, err := s.service.CreateOrder(s.ctx, tc.userUUID, tc.partUUIDs)

			s.Require().NoError(err)
			assert.Equal(s.T(), createdOrder, order)
		})
	}
}

func (s *OrderServiceSuite) TestCreateOrderError() {
	testCases := []struct {
		name          string
		userUUID      string
		partUUIDs     []string
		mockSetup     func()
		expectedError error
	}{
		{
			name:          "Empty part UUIDs",
			userUUID:      "123e4567-e89b-12d3-a456-426614174000",
			partUUIDs:     []string{},
			mockSetup:     func() {},
			expectedError: model.ErrBadRequest,
		},
		{
			name:      "Inventory client error",
			userUUID:  "123e4567-e89b-12d3-a456-426614174000",
			partUUIDs: []string{"part-uuid-1"},
			mockSetup: func() {
				filter := &model.PartsFilter{
					UUIDs: []string{"part-uuid-1"},
				}
				s.inventoryClient.On("ListParts", s.ctx, filter).
					Return(nil, ErrInventoryServiceError).Once()
			},
			expectedError: ErrInventoryServiceError,
		},
		{
			name:      "No parts found",
			userUUID:  "123e4567-e89b-12d3-a456-426614174000",
			partUUIDs: []string{"part-uuid-1"},
			mockSetup: func() {
				filter := &model.PartsFilter{
					UUIDs: []string{"part-uuid-1"},
				}
				s.inventoryClient.On("ListParts", s.ctx, filter).
					Return([]model.Part{}, nil).Once()
			},
			expectedError: model.ErrPartsNotFound,
		},
		{
			name:      "Not all parts found",
			userUUID:  "123e4567-e89b-12d3-a456-426614174000",
			partUUIDs: []string{"part-uuid-1", "part-uuid-2", "part-uuid-3"},
			mockSetup: func() {
				filter := &model.PartsFilter{
					UUIDs: []string{"part-uuid-1", "part-uuid-2", "part-uuid-3"},
				}
				parts := []model.Part{
					{UUID: "part-uuid-1", Name: "Engine", Price: 100.00},
					{UUID: "part-uuid-2", Name: "Fuel", Price: 50.00},
				}
				s.inventoryClient.On("ListParts", s.ctx, filter).
					Return(parts, nil).Once()
			},
			expectedError: model.ErrPartsNotFound,
		},
		{
			name:      "Repository create error",
			userUUID:  "123e4567-e89b-12d3-a456-426614174000",
			partUUIDs: []string{"part-uuid-1"},
			mockSetup: func() {
				filter := &model.PartsFilter{
					UUIDs: []string{"part-uuid-1"},
				}
				parts := []model.Part{
					{UUID: "part-uuid-1", Name: "Engine", Price: 100.00},
				}
				s.inventoryClient.On("ListParts", s.ctx, filter).
					Return(parts, nil).Once()

				expectedOrder := &model.Order{
					UserUUID:    "123e4567-e89b-12d3-a456-426614174000",
					PartUUIDs:   []string{"part-uuid-1"},
					TotalPrice:  100.00,
					OrderStatus: model.OrderStatusPENDINGPAYMENT,
				}
				s.orderRepository.On("CreateOrder", s.ctx, expectedOrder).
					Return(nil, ErrCreateOrderError).Once()
			},
			expectedError: ErrCreateOrderError,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			tc.mockSetup()

			order, err := s.service.CreateOrder(s.ctx, tc.userUUID, tc.partUUIDs)

			assert.ErrorIs(s.T(), err, tc.expectedError)
			assert.Nil(s.T(), order)
		})
	}
}

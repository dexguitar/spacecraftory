package order

import (
	"errors"

	"github.com/stretchr/testify/assert"

	"github.com/dexguitar/spacecraftory/order/internal/model"
)

var ErrUpdateOrder = errors.New("update_order error")

func (s *OrderServiceSuite) TestCancelOrderSuccess() {
	order := &model.Order{
		OrderUUID:       "123e4567-e89b-12d3-a456-426614174000",
		UserUUID:        "123e4567-e89b-12d3-a456-426614174012",
		PartUUIDs:       []string{"123e4567-e89b-12d3-a456-426614174001", "123e4567-e89b-12d3-a456-426614174002"},
		TotalPrice:      100.0,
		OrderStatus:     model.OrderStatusPENDINGPAYMENT,
		TransactionUUID: "123e4567-e89b-12d3-a456-426614174003",
		PaymentMethod:   model.PaymentMethodCARD,
	}

	s.orderRepository.On("GetOrder", s.ctx, order.OrderUUID).
		Return(order, nil).Once()

	s.orderRepository.On("UpdateOrder", s.ctx, order).
		Return(nil).Once()

	err := s.service.CancelOrder(s.ctx, order.OrderUUID)

	s.Require().NoError(err)
	assert.Equal(s.T(), model.OrderStatusCANCELLED, order.OrderStatus)
}

func (s *OrderServiceSuite) TestCancelOrderError() {
	testCases := []struct {
		name          string
		orderUUID     string
		mockSetup     func()
		expectedError error
	}{
		{
			name:      "Order not found",
			orderUUID: "non-existent-uuid",
			mockSetup: func() {
				s.orderRepository.On("GetOrder", s.ctx, "non-existent-uuid").
					Return(nil, model.ErrOrderNotFound).Once()
			},
			expectedError: model.ErrOrderNotFound,
		},
		{
			name:      "Order already paid",
			orderUUID: "123e4567-e89b-12d3-a456-426614174000",
			mockSetup: func() {
				order := &model.Order{
					OrderUUID:   "123e4567-e89b-12d3-a456-426614174000",
					OrderStatus: model.OrderStatusPAID,
				}
				s.orderRepository.On("GetOrder", s.ctx, "123e4567-e89b-12d3-a456-426614174000").
					Return(order, nil).Once()
			},
			expectedError: model.ErrInvalidOrderStatus,
		},
		{
			name:      "Order already cancelled",
			orderUUID: "123e4567-e89b-12d3-a456-426614174000",
			mockSetup: func() {
				order := &model.Order{
					OrderUUID:   "123e4567-e89b-12d3-a456-426614174000",
					OrderStatus: model.OrderStatusCANCELLED,
				}
				s.orderRepository.On("GetOrder", s.ctx, "123e4567-e89b-12d3-a456-426614174000").
					Return(order, nil).Once()
			},
			expectedError: model.ErrInvalidOrderStatus,
		},
		{
			name:      "Repository update error",
			orderUUID: "123e4567-e89b-12d3-a456-426614174000",
			mockSetup: func() {
				order := &model.Order{
					OrderUUID:   "123e4567-e89b-12d3-a456-426614174000",
					OrderStatus: model.OrderStatusPENDINGPAYMENT,
				}
				s.orderRepository.On("GetOrder", s.ctx, "123e4567-e89b-12d3-a456-426614174000").
					Return(order, nil).Once()

				s.orderRepository.On("UpdateOrder", s.ctx, order).
					Return(model.ErrOrderNotFound).Once()
			},
			expectedError: model.ErrOrderNotFound,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			tc.mockSetup()

			err := s.service.CancelOrder(s.ctx, tc.orderUUID)

			assert.Equal(s.T(), err.Error(), tc.expectedError.Error())
		})
	}
}

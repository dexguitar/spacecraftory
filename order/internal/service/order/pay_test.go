package order

import (
	"errors"

	"github.com/stretchr/testify/assert"

	"github.com/dexguitar/spacecraftory/order/internal/model"
)

var (
	ErrPaymentClientError = errors.New("payment client error")
	ErrUpdateOrderError   = errors.New("update order error")
)

func (s *OrderServiceSuite) TestPayOrderSuccess() {
	testCases := []struct {
		name            string
		orderUUID       string
		paymentMethod   model.PaymentMethod
		order           *model.Order
		transactionUUID string
	}{
		{
			name:          "Payment with CARD",
			orderUUID:     "123e4567-e89b-12d3-a456-426614174000",
			paymentMethod: model.PaymentMethodCARD,
			order: &model.Order{
				OrderUUID:   "123e4567-e89b-12d3-a456-426614174000",
				UserUUID:    "123e4567-e89b-12d3-a456-426614174012",
				OrderStatus: model.OrderStatusPENDINGPAYMENT,
			},
			transactionUUID: "txn-card-123",
		},
		{
			name:          "Payment with SBP",
			orderUUID:     "123e4567-e89b-12d3-a456-426614174001",
			paymentMethod: model.PaymentMethodSBP,
			order: &model.Order{
				OrderUUID:   "123e4567-e89b-12d3-a456-426614174001",
				UserUUID:    "123e4567-e89b-12d3-a456-426614174013",
				OrderStatus: model.OrderStatusPENDINGPAYMENT,
			},
			transactionUUID: "txn-sbp-456",
		},
		{
			name:          "Payment with CREDIT_CARD",
			orderUUID:     "123e4567-e89b-12d3-a456-426614174002",
			paymentMethod: model.PaymentMethodCREDIT_CARD,
			order: &model.Order{
				OrderUUID:   "123e4567-e89b-12d3-a456-426614174002",
				UserUUID:    "123e4567-e89b-12d3-a456-426614174014",
				OrderStatus: model.OrderStatusPENDINGPAYMENT,
			},
			transactionUUID: "txn-credit-789",
		},
		{
			name:          "Payment with INVESTOR_MONEY",
			orderUUID:     "123e4567-e89b-12d3-a456-426614174003",
			paymentMethod: model.PaymentMethodINVESTOR_MONEY,
			order: &model.Order{
				OrderUUID:   "123e4567-e89b-12d3-a456-426614174003",
				UserUUID:    "123e4567-e89b-12d3-a456-426614174015",
				OrderStatus: model.OrderStatusPENDINGPAYMENT,
			},
			transactionUUID: "txn-investor-abc",
		},
		{
			name:          "Payment with UNKNOWN method",
			orderUUID:     "123e4567-e89b-12d3-a456-426614174003",
			paymentMethod: model.PaymentMethodUNKNOWN,
			order: &model.Order{
				OrderUUID:   "123e4567-e89b-12d3-a456-426614174003",
				UserUUID:    "123e4567-e89b-12d3-a456-426614174015",
				OrderStatus: model.OrderStatusPENDINGPAYMENT,
			},
			transactionUUID: "txn-unknown-abc",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.orderRepository.On("GetOrder", s.ctx, tc.orderUUID).
				Return(tc.order, nil).Once()

			s.paymentClient.On("PayOrder", s.ctx, tc.orderUUID, tc.order.UserUUID, tc.paymentMethod).
				Return(tc.transactionUUID, nil).Once()

			updatedOrder := &model.Order{
				OrderUUID:       tc.order.OrderUUID,
				UserUUID:        tc.order.UserUUID,
				OrderStatus:     model.OrderStatusPAID,
				TransactionUUID: tc.transactionUUID,
				PaymentMethod:   tc.paymentMethod,
			}

			s.orderRepository.On("UpdateOrder", s.ctx, updatedOrder).
				Return(nil).Once()

			transactionUUID, err := s.service.PayOrder(s.ctx, tc.orderUUID, tc.paymentMethod)

			s.Require().NoError(err)
			assert.Equal(s.T(), tc.transactionUUID, transactionUUID)
		})
	}
}

func (s *OrderServiceSuite) TestPayOrderError() {
	testCases := []struct {
		name          string
		orderUUID     string
		paymentMethod model.PaymentMethod
		mockSetup     func()
		expectedError error
	}{
		{
			name:          "Order not found",
			orderUUID:     "non-existent-uuid",
			paymentMethod: model.PaymentMethodCARD,
			mockSetup: func() {
				s.orderRepository.On("GetOrder", s.ctx, "non-existent-uuid").
					Return(nil, model.ErrOrderNotFound).Once()
			},
			expectedError: model.ErrOrderNotFound,
		},
		{
			name:          "Order already paid",
			orderUUID:     "123e4567-e89b-12d3-a456-426614174000",
			paymentMethod: model.PaymentMethodCARD,
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
			name:          "Order already cancelled",
			orderUUID:     "123e4567-e89b-12d3-a456-426614174000",
			paymentMethod: model.PaymentMethodCARD,
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
			name:          "Payment client error",
			orderUUID:     "123e4567-e89b-12d3-a456-426614174000",
			paymentMethod: model.PaymentMethodCARD,
			mockSetup: func() {
				order := &model.Order{
					OrderUUID:   "123e4567-e89b-12d3-a456-426614174000",
					UserUUID:    "123e4567-e89b-12d3-a456-426614174012",
					OrderStatus: model.OrderStatusPENDINGPAYMENT,
				}
				s.orderRepository.On("GetOrder", s.ctx, "123e4567-e89b-12d3-a456-426614174000").
					Return(order, nil).Once()

				s.paymentClient.On("PayOrder", s.ctx, "123e4567-e89b-12d3-a456-426614174000", order.UserUUID, model.PaymentMethodCARD).
					Return("", ErrPaymentClientError).Once()
			},
			expectedError: model.ErrPaymentFailed,
		},
		{
			name:          "Repository update error",
			orderUUID:     "123e4567-e89b-12d3-a456-426614174000",
			paymentMethod: model.PaymentMethodCARD,
			mockSetup: func() {
				order := &model.Order{
					OrderUUID:   "123e4567-e89b-12d3-a456-426614174000",
					UserUUID:    "123e4567-e89b-12d3-a456-426614174012",
					OrderStatus: model.OrderStatusPENDINGPAYMENT,
				}
				s.orderRepository.On("GetOrder", s.ctx, "123e4567-e89b-12d3-a456-426614174000").
					Return(order, nil).Once()

				transactionUUID := "txn-123"
				s.paymentClient.On("PayOrder", s.ctx, "123e4567-e89b-12d3-a456-426614174000", order.UserUUID, model.PaymentMethodCARD).
					Return(transactionUUID, nil).Once()

				updatedOrder := &model.Order{
					OrderUUID:       order.OrderUUID,
					UserUUID:        order.UserUUID,
					OrderStatus:     model.OrderStatusPAID,
					TransactionUUID: transactionUUID,
					PaymentMethod:   model.PaymentMethodCARD,
				}
				s.orderRepository.On("UpdateOrder", s.ctx, updatedOrder).
					Return(ErrUpdateOrderError).Once()
			},
			expectedError: ErrUpdateOrderError,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			tc.mockSetup()

			transactionUUID, err := s.service.PayOrder(s.ctx, tc.orderUUID, tc.paymentMethod)

			assert.ErrorIs(s.T(), err, tc.expectedError)
			assert.Empty(s.T(), transactionUUID)
		})
	}
}

package v1

import (
	"errors"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/dexguitar/spacecraftory/order/internal/converter"
	"github.com/dexguitar/spacecraftory/order/internal/model"
	orderV1 "github.com/dexguitar/spacecraftory/shared/pkg/openapi/order/v1"
)

func (s *APISuite) TestPayOrderSuccess() {
	testCases := []struct {
		name            string
		orderUUID       uuid.UUID
		paymentMethod   orderV1.PaymentMethod
		transactionUUID string
	}{
		{
			name:            "Payment with CARD",
			orderUUID:       uuid.New(),
			paymentMethod:   orderV1.PaymentMethodCARD,
			transactionUUID: uuid.New().String(),
		},
		{
			name:            "Payment with SBP",
			orderUUID:       uuid.New(),
			paymentMethod:   orderV1.PaymentMethodSBP,
			transactionUUID: uuid.New().String(),
		},
		{
			name:            "Payment with CREDIT_CARD",
			orderUUID:       uuid.New(),
			paymentMethod:   orderV1.PaymentMethodCREDITCARD,
			transactionUUID: uuid.New().String(),
		},
		{
			name:            "Payment with INVESTOR_MONEY",
			orderUUID:       uuid.New(),
			paymentMethod:   orderV1.PaymentMethodINVESTORMONEY,
			transactionUUID: uuid.New().String(),
		},
		{
			name:            "Payment with UNKNOWN",
			orderUUID:       uuid.New(),
			paymentMethod:   orderV1.PaymentMethodUNKNOWN,
			transactionUUID: uuid.New().String(),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			servicePaymentMethod := converter.ToModelPaymentMethod(tc.paymentMethod)

			s.orderService.On("PayOrder", s.ctx, tc.orderUUID.String(), servicePaymentMethod).
				Return(tc.transactionUUID, nil).Once()

			req := &orderV1.PayOrderRequest{
				PaymentMethod: tc.paymentMethod,
			}
			params := orderV1.PayOrderParams{
				OrderUUID: tc.orderUUID,
			}

			resp, err := s.api.PayOrder(s.ctx, req, params)
			s.Require().NoError(err)

			parsedTransactionUUID, err := uuid.Parse(tc.transactionUUID)

			s.Require().NoError(err)
			assert.Equal(s.T(), resp, &orderV1.PayOrderResponse{
				TransactionUUID: parsedTransactionUUID,
			})
		})
	}
}

func (s *APISuite) TestPayOrderError() {
	testCases := []struct {
		name             string
		orderUUID        uuid.UUID
		paymentMethod    orderV1.PaymentMethod
		serviceError     error
		serviceTxnUUID   string
		expectedRespType any
		expectedCode     int
		expectedMessage  string
	}{
		{
			name:             "Order not found",
			orderUUID:        uuid.New(),
			paymentMethod:    orderV1.PaymentMethodCARD,
			serviceError:     model.ErrOrderNotFound,
			expectedRespType: &orderV1.NotFoundError{},
			expectedCode:     404,
			expectedMessage:  "Order not found",
		},
		{
			name:             "Invalid order status - already paid",
			orderUUID:        uuid.New(),
			paymentMethod:    orderV1.PaymentMethodCARD,
			serviceError:     model.ErrInvalidOrderStatus,
			expectedRespType: &orderV1.ConflictError{},
			expectedCode:     409,
			expectedMessage:  "Order has already been paid or cancelled",
		},
		{
			name:             "Service internal error",
			orderUUID:        uuid.New(),
			paymentMethod:    orderV1.PaymentMethodSBP,
			serviceError:     errors.New("payment gateway unavailable"),
			expectedRespType: &orderV1.InternalServerError{},
			expectedCode:     500,
			expectedMessage:  "Failed to process payment",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			servicePaymentMethod := converter.ToModelPaymentMethod(tc.paymentMethod)

			if tc.serviceError != nil {
				s.orderService.On("PayOrder", s.ctx, tc.orderUUID.String(), servicePaymentMethod).
					Return("", tc.serviceError).Once()
			} else if tc.serviceTxnUUID != "" {
				s.orderService.On("PayOrder", s.ctx, tc.orderUUID.String(), servicePaymentMethod).
					Return(tc.serviceTxnUUID, nil).Once()
			}

			req := &orderV1.PayOrderRequest{
				PaymentMethod: tc.paymentMethod,
			}
			params := orderV1.PayOrderParams{
				OrderUUID: tc.orderUUID,
			}

			resp, err := s.api.PayOrder(s.ctx, req, params)

			s.Require().NoError(err)

			switch tc.expectedRespType.(type) {
			case *orderV1.NotFoundError:
				notFoundErr, ok := resp.(*orderV1.NotFoundError)
				s.Require().True(ok, "response should be NotFoundError")
				assert.Equal(s.T(), tc.expectedCode, notFoundErr.Code)
				assert.Equal(s.T(), tc.expectedMessage, notFoundErr.Message)
			case *orderV1.ConflictError:
				conflictErr, ok := resp.(*orderV1.ConflictError)
				s.Require().True(ok, "response should be ConflictError")
				assert.Equal(s.T(), tc.expectedCode, conflictErr.Code)
				assert.Equal(s.T(), tc.expectedMessage, conflictErr.Message)
			case *orderV1.InternalServerError:
				internalErr, ok := resp.(*orderV1.InternalServerError)
				s.Require().True(ok, "response should be BadRequestError")
				assert.Equal(s.T(), tc.expectedCode, internalErr.Code)
				assert.Equal(s.T(), tc.expectedMessage, internalErr.Message)
			default:
				s.Fail("unexpected response type")
			}
		})
	}
}

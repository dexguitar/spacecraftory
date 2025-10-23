package v1

import (
	"errors"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/dexguitar/spacecraftory/payment/internal/model"
	paymentV1 "github.com/dexguitar/spacecraftory/shared/pkg/proto/payment/v1"
)

func (s *APISuite) TestPayOrderSuccess() {
	testCases := []struct {
		name           string
		request        *paymentV1.PayOrderRequest
		serviceTxnUUID string
	}{
		{
			name: "Payment with CARD",
			request: &paymentV1.PayOrderRequest{
				OrderUuid:     "123e4567-e89b-12d3-a456-426614174000",
				UserUuid:      "123e4567-e89b-12d3-a456-426614174012",
				PaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_CARD,
			},
			serviceTxnUUID: "txn-card-123",
		},
		{
			name: "Payment with SBP",
			request: &paymentV1.PayOrderRequest{
				OrderUuid:     "123e4567-e89b-12d3-a456-426614174001",
				UserUuid:      "123e4567-e89b-12d3-a456-426614174013",
				PaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_SBP,
			},
			serviceTxnUUID: "txn-sbp-456",
		},
		{
			name: "Payment with CREDIT_CARD",
			request: &paymentV1.PayOrderRequest{
				OrderUuid:     "123e4567-e89b-12d3-a456-426614174002",
				UserUuid:      "123e4567-e89b-12d3-a456-426614174014",
				PaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD,
			},
			serviceTxnUUID: "txn-credit-789",
		},
		{
			name: "Payment with INVESTOR_MONEY",
			request: &paymentV1.PayOrderRequest{
				OrderUuid:     "123e4567-e89b-12d3-a456-426614174003",
				UserUuid:      "123e4567-e89b-12d3-a456-426614174015",
				PaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY,
			},
			serviceTxnUUID: "txn-investor-abc",
		},
		{
			name: "Payment with UNKNOWN method",
			request: &paymentV1.PayOrderRequest{
				OrderUuid:     "123e4567-e89b-12d3-a456-426614174004",
				UserUuid:      "123e4567-e89b-12d3-a456-426614174016",
				PaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_UNKNOWN_UNSPECIFIED,
			},
			serviceTxnUUID: "txn-unknown-xyz",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			expectedPayment := &model.Payment{
				OrderUUID:     tc.request.OrderUuid,
				UserUUID:      tc.request.UserUuid,
				PaymentMethod: model.PaymentMethodMap[tc.request.PaymentMethod],
			}

			s.paymentService.On("PayOrder", s.ctx, expectedPayment).
				Return(tc.serviceTxnUUID, nil).Once()

			resp, err := s.api.PayOrder(s.ctx, tc.request)

			s.Require().NoError(err)
			assert.Equal(s.T(), tc.serviceTxnUUID, resp.TransactionUuid)
		})
	}
}

func (s *APISuite) TestPayOrderError() {
	testCases := []struct {
		name             string
		request          *paymentV1.PayOrderRequest
		serviceError     error
		expectedCode     codes.Code
		expectedMsgParts []string
	}{
		{
			name: "Invalid order UUID",
			request: &paymentV1.PayOrderRequest{
				OrderUuid:     "invalid-uuid",
				UserUuid:      "123e4567-e89b-12d3-a456-426614174012",
				PaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_CARD,
			},
			expectedCode:     codes.InvalidArgument,
			expectedMsgParts: []string{"Invalid request details"},
		},
		{
			name: "Invalid user UUID",
			request: &paymentV1.PayOrderRequest{
				OrderUuid:     "123e4567-e89b-12d3-a456-426614174000",
				UserUuid:      "invalid-uuid",
				PaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_CARD,
			},
			expectedCode:     codes.InvalidArgument,
			expectedMsgParts: []string{"Invalid request details"},
		},
		{
			name: "Service returns ErrBadRequest",
			request: &paymentV1.PayOrderRequest{
				OrderUuid:     "123e4567-e89b-12d3-a456-426614174000",
				UserUuid:      "123e4567-e89b-12d3-a456-426614174012",
				PaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_CARD,
			},
			serviceError:     model.ErrBadRequest,
			expectedCode:     codes.InvalidArgument,
			expectedMsgParts: []string{"Invalid request details"},
		},
		{
			name: "Service internal error",
			request: &paymentV1.PayOrderRequest{
				OrderUuid:     "123e4567-e89b-12d3-a456-426614174000",
				UserUuid:      "123e4567-e89b-12d3-a456-426614174012",
				PaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_CARD,
			},
			serviceError:     errors.New("database connection failed"),
			expectedCode:     codes.Internal,
			expectedMsgParts: []string{"Internal server error"},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			if tc.serviceError != nil {
				expectedPayment := &model.Payment{
					OrderUUID:     tc.request.OrderUuid,
					UserUUID:      tc.request.UserUuid,
					PaymentMethod: model.PaymentMethodMap[tc.request.PaymentMethod],
				}

				s.paymentService.On("PayOrder", s.ctx, expectedPayment).
					Return("", tc.serviceError).Once()
			}

			resp, err := s.api.PayOrder(s.ctx, tc.request)

			s.Require().Error(err)
			s.Require().Nil(resp)

			st, ok := status.FromError(err)
			s.Require().True(ok)
			assert.Equal(s.T(), tc.expectedCode, st.Code())
			for _, msgPart := range tc.expectedMsgParts {
				assert.Contains(s.T(), st.Message(), msgPart)
			}
		})
	}
}

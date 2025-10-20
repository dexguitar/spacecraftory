package payment

import (
	"errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/dexguitar/spacecraftory/payment/internal/model"
)

var ErrPaymentFailed = errors.New("payment failed")

func (s *ServiceSuite) TestPayOrderSuccess() {
	testCases := []struct {
		name            string
		orderUUID       string
		userUUID        string
		paymentMethod   model.PaymentMethod
		expectedTxnUUID string
	}{
		{
			name:            "Payment with CARD",
			orderUUID:       "123e4567-e89b-12d3-a456-426614174000",
			userUUID:        "123e4567-e89b-12d3-a456-426614174012",
			paymentMethod:   model.PaymentMethodCARD,
			expectedTxnUUID: "txn-card-123",
		},
		{
			name:            "Payment with SBP",
			orderUUID:       "123e4567-e89b-12d3-a456-426614174001",
			userUUID:        "123e4567-e89b-12d3-a456-426614174013",
			paymentMethod:   model.PaymentMethodSBP,
			expectedTxnUUID: "txn-sbp-456",
		},
		{
			name:            "Payment with CREDIT_CARD",
			orderUUID:       "123e4567-e89b-12d3-a456-426614174002",
			userUUID:        "123e4567-e89b-12d3-a456-426614174014",
			paymentMethod:   model.PaymentMethodCREDIT_CARD,
			expectedTxnUUID: "txn-credit-789",
		},
		{
			name:            "Payment with INVESTOR_MONEY",
			orderUUID:       "123e4567-e89b-12d3-a456-426614174003",
			userUUID:        "123e4567-e89b-12d3-a456-426614174015",
			paymentMethod:   model.PaymentMethodINVESTOR_MONEY,
			expectedTxnUUID: "txn-investor-abc",
		},
		{
			name:            "Payment with UNKNOWN method",
			orderUUID:       "123e4567-e89b-12d3-a456-426614174004",
			userUUID:        "123e4567-e89b-12d3-a456-426614174016",
			paymentMethod:   model.PaymentMethodUNKNOWN,
			expectedTxnUUID: "txn-unknown-xyz",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			payment := &model.Payment{
				OrderUUID:     tc.orderUUID,
				UserUUID:      tc.userUUID,
				PaymentMethod: tc.paymentMethod,
			}

			s.paymentRepo.On("PayOrder", s.ctx, payment).
				Return(tc.expectedTxnUUID, nil).Once()

			transactionUUID, err := s.service.PayOrder(s.ctx, payment)

			s.Require().NoError(err)
			assert.Equal(s.T(), tc.expectedTxnUUID, transactionUUID)
		})
	}
}

func (s *ServiceSuite) TestPayOrderError() {
	s.paymentRepo.On("PayOrder", s.ctx, mock.Anything).
		Return("", ErrPaymentFailed).Once()

	payment := &model.Payment{
		OrderUUID:     "123e4567-e89b-12d3-a456-426614174000",
		UserUUID:      "123e4567-e89b-12d3-a456-426614174012",
		PaymentMethod: model.PaymentMethodCARD,
	}

	transactionUUID, err := s.service.PayOrder(s.ctx, payment)

	assert.ErrorIs(s.T(), err, ErrPaymentFailed)
	assert.Empty(s.T(), transactionUUID)
}

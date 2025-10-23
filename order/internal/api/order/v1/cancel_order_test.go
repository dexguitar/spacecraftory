package v1

import (
	"errors"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/dexguitar/spacecraftory/order/internal/model"
	orderV1 "github.com/dexguitar/spacecraftory/shared/pkg/openapi/order/v1"
)

func (s *APISuite) TestCancelOrderSuccess() {
	orderUUID := uuid.New()

	s.orderService.On("CancelOrder", s.ctx, orderUUID.String()).
		Return(nil).Once()

	params := orderV1.CancelOrderParams{
		OrderUUID: orderUUID,
	}

	resp, err := s.api.CancelOrder(s.ctx, params)

	s.Require().NoError(err)
	assert.Equal(s.T(), resp, &orderV1.CancelOrderNoContent{})
}

func (s *APISuite) TestCancelOrderError() {
	testCases := []struct {
		name             string
		orderUUID        uuid.UUID
		serviceError     error
		expectedRespType interface{}
		expectedCode     int
		expectedMessage  string
	}{
		{
			name:             "Order not found",
			orderUUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			serviceError:     model.ErrOrderNotFound,
			expectedRespType: &orderV1.NotFoundError{},
			expectedCode:     404,
			expectedMessage:  "Order not found",
		},
		{
			name:             "Invalid order status - already paid",
			orderUUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174001"),
			serviceError:     model.ErrInvalidOrderStatus,
			expectedRespType: &orderV1.ConflictError{},
			expectedCode:     409,
			expectedMessage:  "Cannot cancel already paid or cancelled order",
		},
		{
			name:             "Internal server error",
			orderUUID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174002"),
			serviceError:     errors.New("database error"),
			expectedRespType: &orderV1.InternalServerError{},
			expectedCode:     500,
			expectedMessage:  "Failed to cancel order",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.orderService.On("CancelOrder", s.ctx, tc.orderUUID.String()).
				Return(tc.serviceError).Once()

			params := orderV1.CancelOrderParams{
				OrderUUID: tc.orderUUID,
			}

			resp, err := s.api.CancelOrder(s.ctx, params)

			s.Require().NoError(err)
			s.Require().NotNil(resp)

			switch expected := tc.expectedRespType.(type) {
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
				s.Require().True(ok, "response should be InternalServerError")
				assert.Equal(s.T(), tc.expectedCode, internalErr.Code)
				assert.Equal(s.T(), tc.expectedMessage, internalErr.Message)

			default:
				s.Fail("unexpected response type", expected)
			}
		})
	}
}

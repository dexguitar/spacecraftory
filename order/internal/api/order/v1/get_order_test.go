package v1

import (
	"errors"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/dexguitar/spacecraftory/order/internal/model"
	orderV1 "github.com/dexguitar/spacecraftory/shared/pkg/openapi/order/v1"
)

func (s *APISuite) TestGetOrderByUUIDSuccess() {
	orderUUID := uuid.New()

	serviceOrder := &model.Order{
		OrderUUID:   orderUUID.String(),
		UserUUID:    uuid.New().String(),
		PartUUIDs:   []string{uuid.New().String(), uuid.New().String()},
		TotalPrice:  250.75,
		OrderStatus: model.OrderStatusPENDINGPAYMENT,
	}

	s.orderService.On("GetOrder", s.ctx, orderUUID.String()).
		Return(serviceOrder, nil).Once()

	params := orderV1.GetOrderByUUIDParams{
		OrderUUID: orderUUID,
	}

	resp, err := s.api.GetOrderByUUID(s.ctx, params)

	s.Require().NoError(err)

	getResp, ok := resp.(*orderV1.OrderDto)
	s.Require().True(ok, "response should be OrderDto")

	assert.Equal(s.T(), orderUUID, getResp.GetOrderUUID())
	assert.Equal(s.T(), serviceOrder.TotalPrice, getResp.GetTotalPrice())
	assert.Equal(s.T(), orderV1.OrderStatusPENDINGPAYMENT, getResp.GetStatus())
}

func (s *APISuite) TestGetOrderByUUIDError() {
	testCases := []struct {
		name             string
		orderUUID        uuid.UUID
		serviceError     error
		expectedRespType any
		expectedCode     int
		expectedMessage  string
	}{
		{
			name:             "Order not found",
			orderUUID:        uuid.New(),
			serviceError:     model.ErrOrderNotFound,
			expectedRespType: &orderV1.NotFoundError{},
			expectedCode:     404,
			expectedMessage:  "Order not found",
		},
		{
			name:             "Service internal error",
			orderUUID:        uuid.New(),
			serviceError:     errors.New("database connection failed"),
			expectedRespType: &orderV1.InternalServerError{},
			expectedCode:     500,
			expectedMessage:  "Failed to get order",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.orderService.On("GetOrder", s.ctx, tc.orderUUID.String()).
				Return(nil, tc.serviceError).Once()

			params := orderV1.GetOrderByUUIDParams{
				OrderUUID: tc.orderUUID,
			}

			resp, err := s.api.GetOrderByUUID(s.ctx, params)

			s.Require().NoError(err)

			switch tc.expectedRespType.(type) {
			case *orderV1.NotFoundError:
				notFoundErr, ok := resp.(*orderV1.NotFoundError)
				s.Require().True(ok, "response should be NotFoundError")
				assert.Equal(s.T(), tc.expectedCode, notFoundErr.Code)
				assert.Equal(s.T(), tc.expectedMessage, notFoundErr.Message)
			case *orderV1.InternalServerError:
				internalErr, ok := resp.(*orderV1.InternalServerError)
				s.Require().True(ok, "response should be InternalServerError")
				assert.Equal(s.T(), tc.expectedCode, internalErr.Code)
				assert.Equal(s.T(), tc.expectedMessage, internalErr.Message)
			default:
				s.Fail("unexpected response type")
			}
		})
	}
}

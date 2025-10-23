package v1

import (
	"errors"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/dexguitar/spacecraftory/order/internal/model"
	orderV1 "github.com/dexguitar/spacecraftory/shared/pkg/openapi/order/v1"
)

func (s *APISuite) TestCreateOrderSuccess() {
	testCases := []struct {
		name          string
		userUUID      uuid.UUID
		partUUIDs     []uuid.UUID
		createdOrder  *model.Order
		expectedPrice float64
	}{
		{
			name:      "Single part order",
			userUUID:  uuid.New(),
			partUUIDs: []uuid.UUID{uuid.New()},
			createdOrder: &model.Order{
				OrderUUID:   uuid.New().String(),
				UserUUID:    uuid.New().String(),
				PartUUIDs:   []string{uuid.New().String()},
				TotalPrice:  100.50,
				OrderStatus: model.OrderStatusPENDINGPAYMENT,
			},
			expectedPrice: 100.50,
		},
		{
			name:      "Multiple parts order",
			userUUID:  uuid.New(),
			partUUIDs: []uuid.UUID{uuid.New(), uuid.New(), uuid.New()},
			createdOrder: &model.Order{
				OrderUUID:   uuid.New().String(),
				UserUUID:    uuid.New().String(),
				PartUUIDs:   []string{uuid.New().String(), uuid.New().String(), uuid.New().String()},
				TotalPrice:  350.75,
				OrderStatus: model.OrderStatusPENDINGPAYMENT,
			},
			expectedPrice: 350.75,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			partUUIDStrings := make([]string, len(tc.partUUIDs))
			for i, partUUID := range tc.partUUIDs {
				partUUIDStrings[i] = partUUID.String()
			}

			s.orderService.On("CreateOrder", s.ctx, tc.userUUID.String(), partUUIDStrings).
				Return(tc.createdOrder, nil).Once()

			req := &orderV1.CreateOrderRequest{
				UserUUID:  tc.userUUID,
				PartUuids: tc.partUUIDs,
			}

			resp, err := s.api.CreateOrder(s.ctx, req)

			s.Require().NoError(err)
			createResp, ok := resp.(*orderV1.CreateOrderResponse)
			s.Require().True(ok, "response should be CreateOrderResponse")
			assert.Equal(s.T(), tc.createdOrder.OrderUUID, createResp.GetOrderUUID().String())
			assert.Equal(s.T(), tc.expectedPrice, createResp.GetTotalPrice())
		})
	}
}

func (s *APISuite) TestCreateOrderError() {
	testCases := []struct {
		name             string
		userUUID         uuid.UUID
		partUUIDs        []uuid.UUID
		serviceError     error
		serviceOrder     *model.Order
		expectedRespType any
		expectedCode     int
		expectedMessage  string
	}{
		{
			name:             "Empty part UUIDs - bad request",
			userUUID:         uuid.New(),
			partUUIDs:        []uuid.UUID{},
			serviceError:     model.ErrBadRequest,
			expectedRespType: &orderV1.BadRequestError{},
			expectedCode:     400,
			expectedMessage:  model.ErrBadRequest.Error(),
		},
		{
			name:             "Parts not found",
			userUUID:         uuid.New(),
			partUUIDs:        []uuid.UUID{uuid.New(), uuid.New()},
			serviceError:     model.ErrPartsNotFound,
			expectedRespType: &orderV1.BadRequestError{},
			expectedCode:     400,
			expectedMessage:  model.ErrPartsNotFound.Error(),
		},
		{
			name:             "Service internal error",
			userUUID:         uuid.New(),
			partUUIDs:        []uuid.UUID{uuid.New()},
			serviceError:     errors.New("database connection failed"),
			expectedRespType: &orderV1.InternalServerError{},
			expectedCode:     500,
			expectedMessage:  "Failed to create order",
		},
		{
			name:      "Invalid order UUID from service",
			userUUID:  uuid.New(),
			partUUIDs: []uuid.UUID{uuid.New()},
			serviceOrder: &model.Order{
				OrderUUID:   "invalid-uuid-format",
				UserUUID:    uuid.New().String(),
				PartUUIDs:   []string{uuid.New().String()},
				TotalPrice:  100.00,
				OrderStatus: model.OrderStatusPENDINGPAYMENT,
			},
			expectedRespType: &orderV1.InternalServerError{},
			expectedCode:     500,
			expectedMessage:  "Failed to parse order UUID",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			partUUIDStrings := make([]string, len(tc.partUUIDs))
			for i, partUUID := range tc.partUUIDs {
				partUUIDStrings[i] = partUUID.String()
			}

			if tc.serviceError != nil {
				s.orderService.On("CreateOrder", s.ctx, tc.userUUID.String(), partUUIDStrings).
					Return(nil, tc.serviceError).Once()
			} else if tc.serviceOrder != nil {
				s.orderService.On("CreateOrder", s.ctx, tc.userUUID.String(), partUUIDStrings).
					Return(tc.serviceOrder, nil).Once()
			}

			req := &orderV1.CreateOrderRequest{
				UserUUID:  tc.userUUID,
				PartUuids: tc.partUUIDs,
			}

			resp, err := s.api.CreateOrder(s.ctx, req)

			s.Require().NoError(err)

			switch tc.expectedRespType.(type) {
			case *orderV1.BadRequestError:
				badRequestErr, ok := resp.(*orderV1.BadRequestError)
				s.Require().True(ok, "response should be BadRequestError")
				assert.Equal(s.T(), tc.expectedCode, badRequestErr.Code)
				assert.Equal(s.T(), tc.expectedMessage, badRequestErr.Message)
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

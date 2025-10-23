package order

import (
	"github.com/stretchr/testify/assert"

	"github.com/dexguitar/spacecraftory/order/internal/model"
)

func (s *OrderServiceSuite) TestGetOrderSuccess() {
	expectedOrder := &model.Order{
		OrderUUID:       "123e4567-e89b-12d3-a456-426614174000",
		UserUUID:        "123e4567-e89b-12d3-a456-426614174012",
		PartUUIDs:       []string{"part-uuid-1", "part-uuid-2"},
		TotalPrice:      200.00,
		OrderStatus:     model.OrderStatusPENDINGPAYMENT,
		TransactionUUID: "",
		PaymentMethod:   "",
	}

	s.orderRepository.On("GetOrder", s.ctx, expectedOrder.OrderUUID).
		Return(expectedOrder, nil).Once()

	order, err := s.service.GetOrder(s.ctx, expectedOrder.OrderUUID)

	s.Require().NoError(err)
	assert.Equal(s.T(), expectedOrder, order)
}

func (s *OrderServiceSuite) TestGetOrderError() {
	invalidOrderUUID := "non-existent-uuid"

	s.orderRepository.On("GetOrder", s.ctx, invalidOrderUUID).
		Return(nil, model.ErrOrderNotFound).Once()

	order, err := s.service.GetOrder(s.ctx, invalidOrderUUID)

	assert.ErrorIs(s.T(), err, model.ErrOrderNotFound)
	assert.Nil(s.T(), order)
}

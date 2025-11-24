package order

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	clientMocks "github.com/dexguitar/spacecraftory/order/internal/client/mocks"
	"github.com/dexguitar/spacecraftory/order/internal/repository/mocks"
	serviceMocks "github.com/dexguitar/spacecraftory/order/internal/service/mocks"
)

type OrderServiceSuite struct {
	suite.Suite
	ctx             context.Context
	orderRepository *mocks.OrderRepository
	inventoryClient *clientMocks.InventoryClient
	paymentClient   *clientMocks.PaymentClient
	producerService *serviceMocks.ProducerService
	service         *service
}

func (s *OrderServiceSuite) SetupTest() {
	s.ctx = context.Background()
	s.orderRepository = mocks.NewOrderRepository(s.T())
	s.inventoryClient = clientMocks.NewInventoryClient(s.T())
	s.paymentClient = clientMocks.NewPaymentClient(s.T())
	s.producerService = serviceMocks.NewProducerService(s.T())
	s.service = NewService(
		s.orderRepository,
		s.inventoryClient,
		s.paymentClient,
		s.producerService,
	)
}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(OrderServiceSuite))
}

package order

import (
	invClient "github.com/dexguitar/spacecraftory/order/internal/client/inventory"
	payClient "github.com/dexguitar/spacecraftory/order/internal/client/payment"
	"github.com/dexguitar/spacecraftory/order/internal/repository"
	def "github.com/dexguitar/spacecraftory/order/internal/service"
)

type service struct {
	orderRepository repository.OrderRepository
	inventoryClient invClient.InventoryClient
	paymentClient   payClient.PaymentClient
	producerService def.ProducerService
}

func NewService(
	orderRepository repository.OrderRepository,
	inventoryClient invClient.InventoryClient,
	paymentClient payClient.PaymentClient,
	producerService def.ProducerService,
) *service {
	return &service{
		orderRepository: orderRepository,
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
		producerService: producerService,
	}
}

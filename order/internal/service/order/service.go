package order

import (
	client "github.com/dexguitar/spacecraftory/order/internal/client"
	"github.com/dexguitar/spacecraftory/order/internal/repository"
	def "github.com/dexguitar/spacecraftory/order/internal/service"
)

type service struct {
	orderRepository repository.OrderRepository
	inventoryClient client.InventoryClient
	paymentClient   client.PaymentClient
	iamClient       client.IAMClient
	producerService def.ProducerService
}

func NewService(
	orderRepository repository.OrderRepository,
	inventoryClient client.InventoryClient,
	paymentClient client.PaymentClient,
	iamClient client.IAMClient,
	producerService def.ProducerService,
) *service {
	return &service{
		orderRepository: orderRepository,
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
		iamClient:       iamClient,
		producerService: producerService,
	}
}

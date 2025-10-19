package order

import (
	"github.com/dexguitar/spacecraftory/order/internal/client"
	"github.com/dexguitar/spacecraftory/order/internal/repository"
)

type service struct {
	orderRepository repository.OrderRepository
	inventoryClient client.InventoryClient
	paymentClient   client.PaymentClient
}

func NewService(
	orderRepository repository.OrderRepository,
	inventoryClient client.InventoryClient,
	paymentClient client.PaymentClient,
) *service {
	return &service{
		orderRepository: orderRepository,
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
	}
}

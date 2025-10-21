package order

import (
	invClient "github.com/dexguitar/spacecraftory/order/internal/client/inventory"
	payClient "github.com/dexguitar/spacecraftory/order/internal/client/payment"
	"github.com/dexguitar/spacecraftory/order/internal/repository"
)

type service struct {
	orderRepository repository.OrderRepository
	inventoryClient invClient.InventoryClient
	paymentClient   payClient.PaymentClient
}

func NewService(
	orderRepository repository.OrderRepository,
	inventoryClient invClient.InventoryClient,
	paymentClient payClient.PaymentClient,
) *service {
	return &service{
		orderRepository: orderRepository,
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
	}
}

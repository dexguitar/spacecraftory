package order

import (
	"context"

	"github.com/google/uuid"

	serviceModel "github.com/dexguitar/spacecraftory/order/internal/model"
	repoConverter "github.com/dexguitar/spacecraftory/order/internal/repository/converter"
)

func (r *orderRepository) CreateOrder(ctx context.Context, order *serviceModel.Order) (*serviceModel.Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	repoOrder := repoConverter.OrderServiceToRepoModel(order)

	newOrderUUID := uuid.New().String()
	repoOrder.OrderUUID = newOrderUUID

	r.orders[newOrderUUID] = repoOrder

	return repoConverter.OrderRepoToServiceModel(repoOrder), nil
}

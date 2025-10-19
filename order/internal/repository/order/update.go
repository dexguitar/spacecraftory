package order

import (
	"context"

	serviceModel "github.com/dexguitar/spacecraftory/order/internal/model"
	repoConverter "github.com/dexguitar/spacecraftory/order/internal/repository/converter"
)

func (r *orderRepository) UpdateOrder(ctx context.Context, order *serviceModel.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.orders[order.OrderUUID]; !ok {
		return serviceModel.ErrOrderNotFound
	}

	repoOrder := repoConverter.OrderServiceToRepoModel(order)
	r.orders[order.OrderUUID] = repoOrder

	return nil
}

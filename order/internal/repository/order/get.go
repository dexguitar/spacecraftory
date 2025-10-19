package order

import (
	"context"

	serviceModel "github.com/dexguitar/spacecraftory/order/internal/model"
	"github.com/dexguitar/spacecraftory/order/internal/repository/converter"
)

func (r *orderRepository) GetOrder(ctx context.Context, orderUUID string) (*serviceModel.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, ok := r.orders[orderUUID]
	if !ok {
		return nil, serviceModel.ErrOrderNotFound
	}

	return converter.OrderRepoToServiceModel(order), nil
}

package order

import (
	"sync"

	repoModel "github.com/dexguitar/spacecraftory/order/internal/repository/model"
)

type orderRepository struct {
	mu     sync.RWMutex
	orders map[string]*repoModel.Order
}

func NewOrderRepository() *orderRepository {
	return &orderRepository{
		orders: make(map[string]*repoModel.Order),
	}
}

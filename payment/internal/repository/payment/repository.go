package payment

import (
	"sync"

	repoModel "github.com/dexguitar/spacecraftory/payment/internal/repository/model"
)

type paymentRepository struct {
	mu   sync.RWMutex
	data map[string]*repoModel.Payment
}

func NewPaymentRepository() *paymentRepository {
	return &paymentRepository{data: make(map[string]*repoModel.Payment)}
}

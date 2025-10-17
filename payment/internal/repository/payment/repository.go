package payment

import (
	"sync"

	"github.com/google/uuid"

	repoModel "github.com/dexguitar/spacecraftory/payment/internal/repository/model"
)

type paymentRepository struct {
	mu   sync.RWMutex
	data map[uuid.UUID]*repoModel.Payment
}

func NewPaymentRepository() *paymentRepository {
	return &paymentRepository{data: make(map[uuid.UUID]*repoModel.Payment)}
}

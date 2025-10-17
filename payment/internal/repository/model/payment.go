package converter

import (
	"github.com/google/uuid"

	"github.com/dexguitar/spacecraftory/payment/internal/model"
)

type Payment struct {
	OrderUUID     uuid.UUID           `db:"order_uuid"`
	UserUUID      uuid.UUID           `db:"user_uuid"`
	PaymentMethod model.PaymentMethod `db:"payment_method"`
}

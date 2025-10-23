package converter

import (
	"github.com/dexguitar/spacecraftory/payment/internal/model"
)

type Payment struct {
	OrderUUID     string              `db:"order_uuid"`
	UserUUID      string              `db:"user_uuid"`
	PaymentMethod model.PaymentMethod `db:"payment_method"`
}

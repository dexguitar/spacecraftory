package model

import (
	"github.com/dexguitar/spacecraftory/order/internal/model"
)

type Order struct {
	OrderUUID       string               `db:"id"`
	UserUUID        string               `db:"user_uuid"`
	TotalPrice      float64              `db:"total_price"`
	Status          model.OrderStatus    `db:"status"`
	TransactionUUID *string              `db:"transaction_uuid"`
	PaymentMethod   *model.PaymentMethod `db:"payment_method"`
}

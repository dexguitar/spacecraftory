package model

import (
	paymentV1 "github.com/dexguitar/spacecraftory/shared/pkg/proto/payment/v1"
)

type PaymentMethod string

type Payment struct {
	OrderUUID     string
	UserUUID      string
	PaymentMethod PaymentMethod
}

const (
	PaymentMethodCARD           PaymentMethod = "CARD"
	PaymentMethodSBP            PaymentMethod = "SBP"
	PaymentMethodCREDIT_CARD    PaymentMethod = "CREDIT_CARD"
	PaymentMethodINVESTOR_MONEY PaymentMethod = "INVESTOR_MONEY"
	PaymentMethodUNKNOWN        PaymentMethod = "UNKNOWN"
)

var PaymentMethodMap = map[paymentV1.PaymentMethod]PaymentMethod{
	paymentV1.PaymentMethod_PAYMENT_METHOD_CARD:                PaymentMethodCARD,
	paymentV1.PaymentMethod_PAYMENT_METHOD_SBP:                 PaymentMethodSBP,
	paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD:         PaymentMethodCREDIT_CARD,
	paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY:      PaymentMethodINVESTOR_MONEY,
	paymentV1.PaymentMethod_PAYMENT_METHOD_UNKNOWN_UNSPECIFIED: PaymentMethodUNKNOWN,
}

package model

type (
	OrderStatus   string
	PaymentMethod string
)

const (
	PaymentMethodUNKNOWN        PaymentMethod = "UNKNOWN"
	PaymentMethodCARD           PaymentMethod = "CARD"
	PaymentMethodSBP            PaymentMethod = "SBP"
	PaymentMethodCREDIT_CARD    PaymentMethod = "CREDIT_CARD"
	PaymentMethodINVESTOR_MONEY PaymentMethod = "INVESTOR_MONEY"
)

const (
	OrderStatusUNKNOWN        OrderStatus = "UNKNOWN"
	OrderStatusPENDINGPAYMENT OrderStatus = "PENDING_PAYMENT"
	OrderStatusPAID           OrderStatus = "PAID"
	OrderStatusCANCELLED      OrderStatus = "CANCELLED"
)

type Order struct {
	OrderUUID       string
	UserUUID        string
	PartUUIDs       []string
	TotalPrice      float64
	OrderStatus     OrderStatus
	TransactionUUID string
	PaymentMethod   PaymentMethod
}

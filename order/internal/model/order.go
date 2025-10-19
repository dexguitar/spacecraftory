package model

type (
	Status        string
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
	StatusUNKNOWN        Status = "UNKNOWN"
	StatusPENDINGPAYMENT Status = "PENDING_PAYMENT"
	StatusPAID           Status = "PAID"
	StatusCANCELLED      Status = "CANCELLED"
)

type Order struct {
	OrderUUID       string
	UserUUID        string
	PartUUIDs       []string
	TotalPrice      float64
	Status          Status
	TransactionUUID string
	PaymentMethod   PaymentMethod
}

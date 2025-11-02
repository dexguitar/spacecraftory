package converter

import (
	serviceModel "github.com/dexguitar/spacecraftory/order/internal/model"
	repoModel "github.com/dexguitar/spacecraftory/order/internal/repository/model"
)

func ToRepoOrder(serviceOrder *serviceModel.Order) *repoModel.Order {
	if serviceOrder == nil {
		return nil
	}

	var transactionUUID *string
	if serviceOrder.TransactionUUID != "" {
		transactionUUID = &serviceOrder.TransactionUUID
	}

	var paymentMethod *serviceModel.PaymentMethod
	if serviceOrder.PaymentMethod != "" {
		paymentMethod = &serviceOrder.PaymentMethod
	}

	return &repoModel.Order{
		OrderUUID:       serviceOrder.OrderUUID,
		UserUUID:        serviceOrder.UserUUID,
		TotalPrice:      serviceOrder.TotalPrice,
		Status:          serviceOrder.OrderStatus,
		TransactionUUID: transactionUUID,
		PaymentMethod:   paymentMethod,
	}
}

func ToModelOrder(repoOrder *repoModel.Order) *serviceModel.Order {
	if repoOrder == nil {
		return nil
	}

	transactionUUID := ""
	if repoOrder.TransactionUUID != nil {
		transactionUUID = *repoOrder.TransactionUUID
	}

	paymentMethod := serviceModel.PaymentMethod("")
	if repoOrder.PaymentMethod != nil {
		paymentMethod = *repoOrder.PaymentMethod
	}

	return &serviceModel.Order{
		OrderUUID:       repoOrder.OrderUUID,
		UserUUID:        repoOrder.UserUUID,
		PartUUIDs:       []string{}, // Will be filled by repository
		TotalPrice:      repoOrder.TotalPrice,
		OrderStatus:     repoOrder.Status,
		TransactionUUID: transactionUUID,
		PaymentMethod:   paymentMethod,
	}
}

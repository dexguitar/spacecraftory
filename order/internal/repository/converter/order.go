package converter

import (
	serviceModel "github.com/dexguitar/spacecraftory/order/internal/model"
	repoModel "github.com/dexguitar/spacecraftory/order/internal/repository/model"
)

func OrderServiceToRepoModel(serviceOrder *serviceModel.Order) *repoModel.Order {
	if serviceOrder == nil {
		return nil
	}

	return &repoModel.Order{
		OrderUUID:       serviceOrder.OrderUUID,
		UserUUID:        serviceOrder.UserUUID,
		PartUUIDs:       serviceOrder.PartUUIDs,
		TotalPrice:      serviceOrder.TotalPrice,
		Status:          serviceOrder.Status,
		TransactionUUID: serviceOrder.TransactionUUID,
		PaymentMethod:   serviceOrder.PaymentMethod,
	}
}

func OrderRepoToServiceModel(repoOrder *repoModel.Order) *serviceModel.Order {
	if repoOrder == nil {
		return nil
	}

	return &serviceModel.Order{
		OrderUUID:       repoOrder.OrderUUID,
		UserUUID:        repoOrder.UserUUID,
		PartUUIDs:       repoOrder.PartUUIDs,
		TotalPrice:      repoOrder.TotalPrice,
		Status:          repoOrder.Status,
		TransactionUUID: repoOrder.TransactionUUID,
		PaymentMethod:   repoOrder.PaymentMethod,
	}
}

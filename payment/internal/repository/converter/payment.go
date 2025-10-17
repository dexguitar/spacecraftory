package converter

import (
	serviceModel "github.com/dexguitar/spacecraftory/payment/internal/model"
	repoModel "github.com/dexguitar/spacecraftory/payment/internal/repository/model"
)

func PaymentInfoToRepoModel(paymentInfo *serviceModel.Payment) repoModel.Payment {
	return repoModel.Payment{
		OrderUUID:     paymentInfo.OrderUUID,
		UserUUID:      paymentInfo.UserUUID,
		PaymentMethod: paymentInfo.PaymentMethod,
	}
}

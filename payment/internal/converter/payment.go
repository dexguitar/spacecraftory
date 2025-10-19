package converter

import (
	"errors"

	"github.com/google/uuid"

	"github.com/dexguitar/spacecraftory/payment/internal/model"
	paymentV1 "github.com/dexguitar/spacecraftory/shared/pkg/proto/payment/v1"
)

func PaymentDtoToServiceModel(paymentDto *paymentV1.PayOrderRequest) (*model.Payment, error) {
	orderUUID, err := uuid.Parse(paymentDto.OrderUuid)
	if err != nil {
		return nil, errors.New("invalid order UUID")
	}
	userUUID, err := uuid.Parse(paymentDto.UserUuid)
	if err != nil {
		return nil, errors.New("invalid user UUID")
	}
	paymentMethod, ok := model.PaymentMethodMap[paymentDto.PaymentMethod]
	if !ok {
		return nil, errors.New("invalid payment method")
	}
	return &model.Payment{
		OrderUUID:     orderUUID.String(),
		UserUUID:      userUUID.String(),
		PaymentMethod: paymentMethod,
	}, nil
}

func PaymentServiceModelToDto(paymentServiceModel *model.Payment) *paymentV1.PayOrderRequest {
	return &paymentV1.PayOrderRequest{
		OrderUuid:     paymentServiceModel.OrderUUID,
		UserUuid:      paymentServiceModel.UserUUID,
		PaymentMethod: paymentMethodToProto(paymentServiceModel.PaymentMethod),
	}
}

func paymentMethodToProto(paymentMethod model.PaymentMethod) paymentV1.PaymentMethod {
	switch paymentMethod {
	case model.PaymentMethodCARD:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_CARD
	case model.PaymentMethodSBP:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_SBP
	case model.PaymentMethodCREDIT_CARD:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case model.PaymentMethodINVESTOR_MONEY:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	default:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_UNKNOWN_UNSPECIFIED
	}
}

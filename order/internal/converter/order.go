package converter

import (
	"github.com/google/uuid"

	"github.com/dexguitar/spacecraftory/order/internal/model"
	orderV1 "github.com/dexguitar/spacecraftory/shared/pkg/openapi/order/v1"
)

func ToDtoOrder(serviceOrder *model.Order) *orderV1.OrderDto {
	if serviceOrder == nil {
		return nil
	}

	orderUUID, err := uuid.Parse(serviceOrder.OrderUUID)
	if err != nil {
		return nil
	}
	userUUID, err := uuid.Parse(serviceOrder.UserUUID)
	if err != nil {
		return nil
	}

	partUUIDs := make([]uuid.UUID, 0, len(serviceOrder.PartUUIDs))
	for _, uuidStr := range serviceOrder.PartUUIDs {
		if partUUID, err := uuid.Parse(uuidStr); err == nil {
			partUUIDs = append(partUUIDs, partUUID)
		}
	}

	dto := &orderV1.OrderDto{
		OrderUUID:  orderUUID,
		UserUUID:   userUUID,
		PartUuids:  partUUIDs,
		TotalPrice: serviceOrder.TotalPrice,
		Status:     ToDtoStatus(serviceOrder.OrderStatus),
	}

	if serviceOrder.TransactionUUID != "" {
		transactionUUID, err := uuid.Parse(serviceOrder.TransactionUUID)
		if err != nil {
			return nil
		}
		dto.TransactionUUID = orderV1.OptNilUUID{
			Value: transactionUUID,
			Set:   true,
			Null:  false,
		}
	}

	if serviceOrder.PaymentMethod != "" && serviceOrder.PaymentMethod != model.PaymentMethodUNKNOWN {
		dto.PaymentMethod = orderV1.OptNilPaymentMethod{
			Value: ToDtoPaymentMethod(serviceOrder.PaymentMethod),
			Set:   true,
			Null:  false,
		}
	}

	return dto
}

func ToDtoStatus(status model.OrderStatus) orderV1.OrderStatus {
	switch status {
	case model.OrderStatusPENDINGPAYMENT:
		return orderV1.OrderStatusPENDINGPAYMENT
	case model.OrderStatusPAID:
		return orderV1.OrderStatusPAID
	case model.OrderStatusCANCELLED:
		return orderV1.OrderStatusCANCELLED
	case model.OrderStatusASSEMBLED:
		return orderV1.OrderStatusASSEMBLED
	default:
		return orderV1.OrderStatusUNKNOWN
	}
}

func ToModelPaymentMethod(apiMethod orderV1.PaymentMethod) model.PaymentMethod {
	switch apiMethod {
	case orderV1.PaymentMethodCARD:
		return model.PaymentMethodCARD
	case orderV1.PaymentMethodSBP:
		return model.PaymentMethodSBP
	case orderV1.PaymentMethodCREDITCARD:
		return model.PaymentMethodCREDIT_CARD
	case orderV1.PaymentMethodINVESTORMONEY:
		return model.PaymentMethodINVESTOR_MONEY
	default:
		return model.PaymentMethodUNKNOWN
	}
}

func ToDtoPaymentMethod(serviceMethod model.PaymentMethod) orderV1.PaymentMethod {
	switch serviceMethod {
	case model.PaymentMethodCARD:
		return orderV1.PaymentMethodCARD
	case model.PaymentMethodSBP:
		return orderV1.PaymentMethodSBP
	case model.PaymentMethodCREDIT_CARD:
		return orderV1.PaymentMethodCREDITCARD
	case model.PaymentMethodINVESTOR_MONEY:
		return orderV1.PaymentMethodINVESTORMONEY
	default:
		return orderV1.PaymentMethodUNKNOWN
	}
}

func ToProtoTransactionUUID(transactionUUID string) uuid.UUID {
	txUUID, err := uuid.Parse(transactionUUID)
	if err != nil {
		return uuid.Nil
	}
	return txUUID
}

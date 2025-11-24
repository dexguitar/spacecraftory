package kafka

import "github.com/dexguitar/spacecraftory/order/internal/model"

type OrderAssembledDecoder interface {
	Decode(data []byte) (model.OrderAssembledEvent, error)
}

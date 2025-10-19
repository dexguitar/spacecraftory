package converter

import (
	"github.com/dexguitar/spacecraftory/order/internal/model"
	inventoryV1 "github.com/dexguitar/spacecraftory/shared/pkg/proto/inventory/v1"
)

func PartProtoToServiceModel(protoPart *inventoryV1.Part) model.Part {
	return model.Part{
		UUID:  protoPart.GetUuid(),
		Name:  protoPart.GetName(),
		Price: protoPart.GetPrice(),
	}
}

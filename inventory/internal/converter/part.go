package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/dexguitar/spacecraftory/inventory/internal/model"
	inventoryV1 "github.com/dexguitar/spacecraftory/shared/pkg/proto/inventory/v1"
)

func ToModelPart(protoPart *inventoryV1.Part) *model.Part {
	if protoPart == nil {
		return nil
	}

	return &model.Part{
		UUID:          protoPart.GetUuid(),
		Name:          protoPart.GetName(),
		Description:   protoPart.GetDescription(),
		Price:         protoPart.GetPrice(),
		StockQuantity: protoPart.GetStockQuantity(),
		Category:      ToModelCategory(protoPart.GetCategory()),
		Dimensions:    ToModelDimensions(protoPart.GetDimensions()),
		Manufacturer:  ToModelManufacturer(protoPart.GetManufacturer()),
		Tags:          protoPart.GetTags(),
		CreatedAt:     protoPart.GetCreatedAt().AsTime(),
		UpdatedAt:     protoPart.GetUpdatedAt().AsTime(),
	}
}

func ToProtoPart(servicePart *model.Part) *inventoryV1.Part {
	if servicePart == nil {
		return nil
	}

	return &inventoryV1.Part{
		Uuid:          servicePart.UUID,
		Name:          servicePart.Name,
		Description:   servicePart.Description,
		Price:         servicePart.Price,
		StockQuantity: servicePart.StockQuantity,
		Category:      ToProtoCategory(servicePart.Category),
		Dimensions:    ToProtoDimensions(servicePart.Dimensions),
		Manufacturer:  ToProtoManufacturer(servicePart.Manufacturer),
		Tags:          servicePart.Tags,
		CreatedAt:     timestamppb.New(servicePart.CreatedAt),
		UpdatedAt:     timestamppb.New(servicePart.UpdatedAt),
	}
}

func ToProtoParts(serviceParts []*model.Part) []*inventoryV1.Part {
	if serviceParts == nil {
		return nil
	}

	protoParts := make([]*inventoryV1.Part, 0, len(serviceParts))
	for _, part := range serviceParts {
		protoParts = append(protoParts, ToProtoPart(part))
	}
	return protoParts
}

func ToModelPartsFilter(protoFilter *inventoryV1.PartsFilter) *model.PartsFilter {
	if protoFilter == nil {
		return nil
	}

	categories := make([]model.Category, 0, len(protoFilter.GetCategories()))
	for _, cat := range protoFilter.GetCategories() {
		categories = append(categories, ToModelCategory(cat))
	}

	return &model.PartsFilter{
		UUIDs:                 protoFilter.GetUuids(),
		Names:                 protoFilter.GetNames(),
		Categories:            categories,
		ManufacturerCountries: protoFilter.GetManufacturerCountries(),
		Tags:                  protoFilter.GetTags(),
	}
}

func ToModelCategory(protoCategory inventoryV1.Category) model.Category {
	if category, ok := model.CategoryMap[protoCategory]; ok {
		return category
	}
	return model.CategoryUnknown
}

func ToProtoCategory(serviceCategory model.Category) inventoryV1.Category {
	switch serviceCategory {
	case model.CategoryEngine:
		return inventoryV1.Category_CATEGORY_ENGINE
	case model.CategoryFuel:
		return inventoryV1.Category_CATEGORY_FUEL
	case model.CategoryPorthole:
		return inventoryV1.Category_CATEGORY_PORTHOLE
	case model.CategoryWing:
		return inventoryV1.Category_CATEGORY_WING
	default:
		return inventoryV1.Category_CATEGORY_UNKNOWN_UNSPECIFIED
	}
}

func ToModelDimensions(protoDims *inventoryV1.Dimensions) *model.Dimensions {
	if protoDims == nil {
		return nil
	}
	return &model.Dimensions{
		Length: protoDims.GetLength(),
		Width:  protoDims.GetWidth(),
		Height: protoDims.GetHeight(),
		Weight: protoDims.GetWeight(),
	}
}

func ToProtoDimensions(serviceDims *model.Dimensions) *inventoryV1.Dimensions {
	if serviceDims == nil {
		return nil
	}
	return &inventoryV1.Dimensions{
		Length: serviceDims.Length,
		Width:  serviceDims.Width,
		Height: serviceDims.Height,
		Weight: serviceDims.Weight,
	}
}

func ToModelManufacturer(protoMan *inventoryV1.Manufacturer) *model.Manufacturer {
	if protoMan == nil {
		return nil
	}
	return &model.Manufacturer{
		Name:    protoMan.GetName(),
		Country: protoMan.GetCountry(),
		Website: protoMan.GetWebsite(),
	}
}

func ToProtoManufacturer(serviceMan *model.Manufacturer) *inventoryV1.Manufacturer {
	if serviceMan == nil {
		return nil
	}
	return &inventoryV1.Manufacturer{
		Name:    serviceMan.Name,
		Country: serviceMan.Country,
		Website: serviceMan.Website,
	}
}

package converter

import (
	"github.com/dexguitar/spacecraftory/inventory/internal/model"
	inventoryV1 "github.com/dexguitar/spacecraftory/shared/pkg/proto/inventory/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func PartProtoToServiceModel(protoPart *inventoryV1.Part) *model.Part {
	if protoPart == nil {
		return nil
	}

	return &model.Part{
		UUID:          protoPart.GetUuid(),
		Name:          protoPart.GetName(),
		Description:   protoPart.GetDescription(),
		Price:         protoPart.GetPrice(),
		StockQuantity: protoPart.GetStockQuantity(),
		Category:      categoryProtoToService(protoPart.GetCategory()),
		Dimensions:    dimensionsProtoToService(protoPart.GetDimensions()),
		Manufacturer:  manufacturerProtoToService(protoPart.GetManufacturer()),
		Tags:          protoPart.GetTags(),
		CreatedAt:     protoPart.GetCreatedAt().AsTime(),
		UpdatedAt:     protoPart.GetUpdatedAt().AsTime(),
	}
}

func PartServiceModelToProto(servicePart *model.Part) *inventoryV1.Part {
	if servicePart == nil {
		return nil
	}

	return &inventoryV1.Part{
		Uuid:          servicePart.UUID,
		Name:          servicePart.Name,
		Description:   servicePart.Description,
		Price:         servicePart.Price,
		StockQuantity: servicePart.StockQuantity,
		Category:      categoryServiceToProto(servicePart.Category),
		Dimensions:    dimensionsServiceToProto(servicePart.Dimensions),
		Manufacturer:  manufacturerServiceToProto(servicePart.Manufacturer),
		Tags:          servicePart.Tags,
		CreatedAt:     timestamppb.New(servicePart.CreatedAt),
		UpdatedAt:     timestamppb.New(servicePart.UpdatedAt),
	}
}

func PartsFilterProtoToServiceModel(protoFilter *inventoryV1.PartsFilter) *model.PartsFilter {
	if protoFilter == nil {
		return nil
	}

	categories := make([]model.Category, 0, len(protoFilter.GetCategories()))
	for _, cat := range protoFilter.GetCategories() {
		categories = append(categories, categoryProtoToService(cat))
	}

	return &model.PartsFilter{
		UUIDs:                 protoFilter.GetUuids(),
		Names:                 protoFilter.GetNames(),
		Categories:            categories,
		ManufacturerCountries: protoFilter.GetManufacturerCountries(),
		Tags:                  protoFilter.GetTags(),
	}
}


func categoryProtoToService(protoCategory inventoryV1.Category) model.Category {
	if category, ok := model.CategoryMap[protoCategory]; ok {
		return category
	}
	return model.CategoryUnknown
}

func categoryServiceToProto(serviceCategory model.Category) inventoryV1.Category {
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

func dimensionsProtoToService(protoDims *inventoryV1.Dimensions) *model.Dimensions {
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

func dimensionsServiceToProto(serviceDims *model.Dimensions) *inventoryV1.Dimensions {
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

func manufacturerProtoToService(protoMan *inventoryV1.Manufacturer) *model.Manufacturer {
	if protoMan == nil {
		return nil
	}
	return &model.Manufacturer{
		Name:    protoMan.GetName(),
		Country: protoMan.GetCountry(),
		Website: protoMan.GetWebsite(),
	}
}

func manufacturerServiceToProto(serviceMan *model.Manufacturer) *inventoryV1.Manufacturer {
	if serviceMan == nil {
		return nil
	}
	return &inventoryV1.Manufacturer{
		Name:    serviceMan.Name,
		Country: serviceMan.Country,
		Website: serviceMan.Website,
	}
}

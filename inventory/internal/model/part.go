package model

import (
	"time"

	inventoryV1 "github.com/dexguitar/spacecraftory/shared/pkg/proto/inventory/v1"
)

type Category string

const (
	CategoryEngine   Category = "ENGINE"
	CategoryFuel     Category = "FUEL"
	CategoryPorthole Category = "PORTHOLE"
	CategoryWing     Category = "WING"
	CategoryUnknown  Category = "UNKNOWN"
)

var CategoryMap = map[inventoryV1.Category]Category{
	inventoryV1.Category_CATEGORY_ENGINE:              CategoryEngine,
	inventoryV1.Category_CATEGORY_FUEL:                CategoryFuel,
	inventoryV1.Category_CATEGORY_PORTHOLE:            CategoryPorthole,
	inventoryV1.Category_CATEGORY_WING:                CategoryWing,
	inventoryV1.Category_CATEGORY_UNKNOWN_UNSPECIFIED: CategoryUnknown,
}

type Dimensions struct {
	Length float64
	Width  float64
	Height float64
	Weight float64
}

type Manufacturer struct {
	Name    string
	Country string
	Website string
}

type Part struct {
	UUID          string
	Name          string
	Description   string
	Price         float64
	StockQuantity int64
	Category      Category
	Dimensions    *Dimensions
	Manufacturer  *Manufacturer
	Tags          []string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type PartsFilter struct {
	UUIDs                 []string
	Names                 []string
	Categories            []Category
	ManufacturerCountries []string
	Tags                  []string
}

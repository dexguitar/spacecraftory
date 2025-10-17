package model

import (
	"time"

	"github.com/dexguitar/spacecraftory/inventory/internal/model"
)

type Dimensions struct {
	Length float64 `db:"length"`
	Width  float64 `db:"width"`
	Height float64 `db:"height"`
	Weight float64 `db:"weight"`
}

type Manufacturer struct {
	Name    string `db:"name"`
	Country string `db:"country"`
	Website string `db:"website"`
}

type Part struct {
	UUID          string         `db:"uuid"`
	Name          string         `db:"name"`
	Description   string         `db:"description"`
	Price         float64        `db:"price"`
	StockQuantity int64          `db:"stock_quantity"`
	Category      model.Category `db:"category"`
	Dimensions    *Dimensions    `db:"dimensions"`
	Manufacturer  *Manufacturer  `db:"manufacturer"`
	Tags          []string       `db:"tags"`
	CreatedAt     time.Time      `db:"created_at"`
	UpdatedAt     time.Time      `db:"updated_at"`
}

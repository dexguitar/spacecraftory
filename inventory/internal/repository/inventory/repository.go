package inventory

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/dexguitar/spacecraftory/inventory/internal/model"
	repoModel "github.com/dexguitar/spacecraftory/inventory/internal/repository/model"
)

type inventoryRepository struct {
	db *mongo.Database
}

func NewInventoryRepository(db *mongo.Database) *inventoryRepository {
	repo := &inventoryRepository{
		db: db,
	}

	repo.initParts()

	return repo
}

func (r *inventoryRepository) initParts() {
	now := time.Now()

	count, err := r.db.Collection("parts").CountDocuments(context.Background(), bson.M{})
	if err != nil {
		log.Printf("failed to count parts: %v\n", err)
		return
	}
	if count > 0 {
		log.Printf("✅ parts already initialized")
		return
	}

	mockParts := map[string]*repoModel.Part{
		"123e4567-e89b-12d3-a456-426614174000": {
			UUID:          "123e4567-e89b-12d3-a456-426614174000",
			Name:          "Quantum Drive Engine",
			Description:   "High-efficiency quantum propulsion engine for interstellar travel",
			Price:         150000.00,
			StockQuantity: 5,
			Category:      model.CategoryEngine,
			Dimensions: &repoModel.Dimensions{
				Length: 3.5,
				Width:  2.0,
				Height: 2.5,
				Weight: 500.0,
			},
			Manufacturer: &repoModel.Manufacturer{
				Name:    "SpaceTech Industries",
				Country: "USA",
				Website: "https://spacetech.example.com",
			},
			Tags:      []string{"quantum", "propulsion", "interstellar"},
			CreatedAt: now,
			UpdatedAt: now,
		},
		"123e4567-e89b-12d3-a456-426614174001": {
			UUID:          "123e4567-e89b-12d3-a456-426614174001",
			Name:          "Fusion Fuel Cell",
			Description:   "Advanced fusion-based fuel cell for long-duration missions",
			Price:         75000.00,
			StockQuantity: 12,
			Category:      model.CategoryFuel,
			Dimensions: &repoModel.Dimensions{
				Length: 1.2,
				Width:  0.8,
				Height: 1.0,
				Weight: 100.0,
			},
			Manufacturer: &repoModel.Manufacturer{
				Name:    "Energy Solutions Corp",
				Country: "Germany",
				Website: "https://energysolutions.example.com",
			},
			Tags:      []string{"fusion", "fuel", "efficient"},
			CreatedAt: now,
			UpdatedAt: now,
		},
		"123e4567-e89b-12d3-a456-426614174002": {
			UUID:          "123e4567-e89b-12d3-a456-426614174002",
			Name:          "Reinforced Porthole",
			Description:   "Triple-layered reinforced viewing porthole for crew observation",
			Price:         25000.00,
			StockQuantity: 20,
			Category:      model.CategoryPorthole,
			Dimensions: &repoModel.Dimensions{
				Length: 1.0,
				Width:  1.0,
				Height: 0.3,
				Weight: 50.0,
			},
			Manufacturer: &repoModel.Manufacturer{
				Name:    "ViewTech Manufacturing",
				Country: "Japan",
				Website: "https://viewtech.example.com",
			},
			Tags:      []string{"observation", "reinforced", "safety"},
			CreatedAt: now,
			UpdatedAt: now,
		},
		"123e4567-e89b-12d3-a456-426614174003": {
			UUID:          "123e4567-e89b-12d3-a456-426614174003",
			Name:          "Aerodynamic Wing Panel",
			Description:   "Carbon-fiber composite wing panel for atmospheric re-entry",
			Price:         45000.00,
			StockQuantity: 8,
			Category:      model.CategoryWing,
			Dimensions: &repoModel.Dimensions{
				Length: 5.0,
				Width:  2.5,
				Height: 0.5,
				Weight: 200.0,
			},
			Manufacturer: &repoModel.Manufacturer{
				Name:    "AeroDynamics Ltd",
				Country: "UK",
				Website: "https://aerodynamics.example.com",
			},
			Tags:      []string{"aerodynamic", "reentry", "composite"},
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	_, err = r.db.Collection("parts").InsertMany(context.Background(), []any{mockParts["123e4567-e89b-12d3-a456-426614174000"], mockParts["123e4567-e89b-12d3-a456-426614174001"], mockParts["123e4567-e89b-12d3-a456-426614174002"], mockParts["123e4567-e89b-12d3-a456-426614174003"]})
	if err != nil {
		log.Printf("failed to insert parts: %v\n", err)
		return
	}

	log.Printf("✅ initialized 4 mock parts")
}

package integration

import (
	"context"
	"os"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	inventoryV1 "github.com/dexguitar/spacecraftory/shared/pkg/proto/inventory/v1"
)

// InsertTestPart inserts a test part with given data
func (env *TestEnvironment) InsertTestPart(ctx context.Context, name, description string, price float64, category inventoryV1.Category) (string, error) {
	partUUID := gofakeit.UUID()
	now := time.Now()

	categoryStr := "UNKNOWN"
	switch category {
	case inventoryV1.Category_CATEGORY_ENGINE:
		categoryStr = "ENGINE"
	case inventoryV1.Category_CATEGORY_FUEL:
		categoryStr = "FUEL"
	case inventoryV1.Category_CATEGORY_PORTHOLE:
		categoryStr = "PORTHOLE"
	case inventoryV1.Category_CATEGORY_WING:
		categoryStr = "WING"
	}

	partDoc := bson.M{
		"uuid":           partUUID,
		"name":           name,
		"description":    description,
		"price":          price,
		"stock_quantity": int64(gofakeit.Number(10, 100)),
		"category":       categoryStr,
		"dimensions": bson.M{
			"length": gofakeit.Float64Range(1, 100),
			"width":  gofakeit.Float64Range(1, 100),
			"height": gofakeit.Float64Range(1, 100),
			"weight": gofakeit.Float64Range(1, 1000),
		},
		"manufacturer": bson.M{
			"name":    gofakeit.Company(),
			"country": gofakeit.Country(),
			"website": gofakeit.URL(),
		},
		"tags":       []string{"test", "e2e"},
		"created_at": primitive.NewDateTimeFromTime(now),
		"updated_at": primitive.NewDateTimeFromTime(now),
	}

	databaseName := os.Getenv("MONGO_DATABASE")
	if databaseName == "" {
		databaseName = "inventory-service"
	}

	_, err := env.Mongo.Client().Database(databaseName).Collection(partsCollectionName).InsertOne(ctx, partDoc)
	if err != nil {
		return "", err
	}

	return partUUID, nil
}

// ClearPartsCollection deletes all records from the parts collection
func (env *TestEnvironment) ClearPartsCollection(ctx context.Context) error {
	databaseName := os.Getenv("MONGO_DATABASE")
	if databaseName == "" {
		databaseName = "inventory-service"
	}

	_, err := env.Mongo.Client().Database(databaseName).Collection(partsCollectionName).DeleteMany(ctx, bson.M{})
	if err != nil {
		return err
	}

	return nil
}

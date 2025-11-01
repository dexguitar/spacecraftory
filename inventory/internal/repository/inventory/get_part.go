package inventory

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/dexguitar/spacecraftory/inventory/internal/model"
	repoConverter "github.com/dexguitar/spacecraftory/inventory/internal/repository/converter"
	repoModel "github.com/dexguitar/spacecraftory/inventory/internal/repository/model"
)

func (r *inventoryRepository) GetPart(ctx context.Context, uuid string) (*model.Part, error) {
	var part repoModel.Part
	err := r.db.Collection("parts").FindOne(ctx, bson.M{"uuid": uuid}).Decode(&part)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, model.ErrPartNotFound
		}

		return nil, err
	}

	return repoConverter.ToModelPart(&part), nil
}

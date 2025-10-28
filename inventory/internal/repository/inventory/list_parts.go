package inventory

import (
	"context"
	"log"
	"slices"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/dexguitar/spacecraftory/inventory/internal/model"
	repoConverter "github.com/dexguitar/spacecraftory/inventory/internal/repository/converter"
	repoModel "github.com/dexguitar/spacecraftory/inventory/internal/repository/model"
)

func (r *inventoryRepository) ListParts(ctx context.Context, filter *model.PartsFilter) ([]*model.Part, error) {
	// If no filter, return all parts
	if filter == nil {
		return r.getAllParts(ctx)
	}

	// UUID filter priority
	if len(filter.UUIDs) > 0 {
		return r.getPartsByUUIDs(ctx, filter.UUIDs)
	}

	// For other filters, fetch all and filter in memory
	return r.filterParts(ctx, filter)
}

func (r *inventoryRepository) getAllParts(ctx context.Context) ([]*model.Part, error) {
	serviceParts := make([]*model.Part, 0)
	cursor, err := r.db.Collection("parts").Find(ctx, bson.M{})
	if err != nil {
		log.Printf("failed to list parts: %v\n", err)
		return nil, err
	}
	defer func() {
		if err := cursor.Close(ctx); err != nil {
			log.Printf("failed to close cursor: %v\n", err)
		}
	}()

	for cursor.Next(ctx) {
		var part repoModel.Part
		if err := cursor.Decode(&part); err != nil {
			log.Printf("failed to decode part: %v\n", err)
			continue
		}
		serviceParts = append(serviceParts, repoConverter.ToModelPart(&part))
	}

	if err := cursor.Err(); err != nil {
		log.Printf("cursor error: %v\n", err)
		return nil, err
	}

	return serviceParts, nil
}

func (r *inventoryRepository) getPartsByUUIDs(ctx context.Context, uuids []string) ([]*model.Part, error) {
	serviceParts := make([]*model.Part, 0, len(uuids))

	// Use MongoDB $in operator for efficient batch query
	filter := bson.M{"uuid": bson.M{"$in": uuids}}
	cursor, err := r.db.Collection("parts").Find(ctx, filter)
	if err != nil {
		log.Printf("failed to find parts by UUIDs: %v\n", err)
		return nil, err
	}
	defer func() {
		if err := cursor.Close(ctx); err != nil {
			log.Printf("failed to close cursor: %v\n", err)
		}
	}()

	for cursor.Next(ctx) {
		var part repoModel.Part
		if err := cursor.Decode(&part); err != nil {
			log.Printf("failed to decode part: %v\n", err)
			continue
		}
		serviceParts = append(serviceParts, repoConverter.ToModelPart(&part))
	}

	if err := cursor.Err(); err != nil {
		log.Printf("cursor error: %v\n", err)
		return nil, err
	}

	return serviceParts, nil
}

func (r *inventoryRepository) filterParts(ctx context.Context, filter *model.PartsFilter) ([]*model.Part, error) {
	// Fetch all parts from database
	allParts, err := r.getAllParts(ctx)
	if err != nil {
		return nil, err
	}

	// Filter in memory
	serviceParts := make([]*model.Part, 0)
	for _, servicePart := range allParts {
		if r.partMatchesFilter(servicePart, filter) {
			serviceParts = append(serviceParts, servicePart)
		}
	}
	return serviceParts, nil
}

func (r *inventoryRepository) partMatchesFilter(part *model.Part, filter *model.PartsFilter) bool {
	return r.matchesNameFilter(part, filter.Names) &&
		r.matchesCategoryFilter(part, filter.Categories) &&
		r.matchesCountryFilter(part, filter.ManufacturerCountries) &&
		r.matchesTagFilter(part, filter.Tags)
}

func (r *inventoryRepository) matchesNameFilter(part *model.Part, names []string) bool {
	if len(names) == 0 {
		return true
	}
	return slices.Contains(names, part.Name)
}

func (r *inventoryRepository) matchesCategoryFilter(part *model.Part, categories []model.Category) bool {
	if len(categories) == 0 {
		return true
	}
	return slices.Contains(categories, part.Category)
}

func (r *inventoryRepository) matchesCountryFilter(part *model.Part, countries []string) bool {
	if len(countries) == 0 {
		return true
	}
	if part.Manufacturer == nil {
		return false
	}
	return slices.Contains(countries, part.Manufacturer.Country)
}

func (r *inventoryRepository) matchesTagFilter(part *model.Part, filterTags []string) bool {
	if len(filterTags) == 0 {
		return true
	}
	return slices.ContainsFunc(part.Tags, func(tag string) bool {
		return slices.Contains(filterTags, tag)
	})
}

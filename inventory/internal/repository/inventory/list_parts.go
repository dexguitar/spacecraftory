package inventory

import (
	"context"
	"slices"

	"github.com/dexguitar/spacecraftory/inventory/internal/model"
	repoConverter "github.com/dexguitar/spacecraftory/inventory/internal/repository/converter"
)

func (r *inventoryRepository) ListParts(ctx context.Context, filter *model.PartsFilter) ([]*model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()


	if filter == nil {
		return r.getAllParts(), nil
	}


	if len(filter.UUIDs) > 0 {
		return r.getPartsByUUIDs(filter.UUIDs), nil
	}


	return r.filterParts(filter), nil
}

func (r *inventoryRepository) getAllParts() []*model.Part {
	serviceParts := make([]*model.Part, 0, len(r.parts))
	for _, repoPart := range r.parts {
		serviceParts = append(serviceParts, repoConverter.PartRepoToServiceModel(repoPart))
	}
	return serviceParts
}

func (r *inventoryRepository) getPartsByUUIDs(uuids []string) []*model.Part {
	serviceParts := make([]*model.Part, 0, len(uuids))
	for _, uuid := range uuids {
		if repoPart, ok := r.parts[uuid]; ok {
			serviceParts = append(serviceParts, repoConverter.PartRepoToServiceModel(repoPart))
		}
	}
	return serviceParts
}

func (r *inventoryRepository) filterParts(filter *model.PartsFilter) []*model.Part {
	serviceParts := make([]*model.Part, 0)
	for _, repoPart := range r.parts {
		servicePart := repoConverter.PartRepoToServiceModel(repoPart)
		if r.partMatchesFilter(servicePart, filter) {
			serviceParts = append(serviceParts, servicePart)
		}
	}
	return serviceParts
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

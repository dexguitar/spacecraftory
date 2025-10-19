package converter

import (
	serviceModel "github.com/dexguitar/spacecraftory/inventory/internal/model"
	repoModel "github.com/dexguitar/spacecraftory/inventory/internal/repository/model"
)

func PartServiceToRepoModel(servicePart *serviceModel.Part) *repoModel.Part {
	if servicePart == nil {
		return nil
	}

	return &repoModel.Part{
		UUID:          servicePart.UUID,
		Name:          servicePart.Name,
		Description:   servicePart.Description,
		Price:         servicePart.Price,
		StockQuantity: servicePart.StockQuantity,
		Category:      servicePart.Category,
		Dimensions:    dimensionsServiceToRepo(servicePart.Dimensions),
		Manufacturer:  manufacturerServiceToRepo(servicePart.Manufacturer),
		Tags:          servicePart.Tags,
		CreatedAt:     servicePart.CreatedAt,
		UpdatedAt:     servicePart.UpdatedAt,
	}
}

func PartRepoToServiceModel(repoPart *repoModel.Part) *serviceModel.Part {
	if repoPart == nil {
		return nil
	}

	return &serviceModel.Part{
		UUID:          repoPart.UUID,
		Name:          repoPart.Name,
		Description:   repoPart.Description,
		Price:         repoPart.Price,
		StockQuantity: repoPart.StockQuantity,
		Category:      repoPart.Category,
		Dimensions:    dimensionsRepoToService(repoPart.Dimensions),
		Manufacturer:  manufacturerRepoToService(repoPart.Manufacturer),
		Tags:          repoPart.Tags,
		CreatedAt:     repoPart.CreatedAt,
		UpdatedAt:     repoPart.UpdatedAt,
	}
}

func dimensionsServiceToRepo(serviceDims *serviceModel.Dimensions) *repoModel.Dimensions {
	if serviceDims == nil {
		return nil
	}
	return &repoModel.Dimensions{
		Length: serviceDims.Length,
		Width:  serviceDims.Width,
		Height: serviceDims.Height,
		Weight: serviceDims.Weight,
	}
}

func dimensionsRepoToService(repoDims *repoModel.Dimensions) *serviceModel.Dimensions {
	if repoDims == nil {
		return nil
	}
	return &serviceModel.Dimensions{
		Length: repoDims.Length,
		Width:  repoDims.Width,
		Height: repoDims.Height,
		Weight: repoDims.Weight,
	}
}

func manufacturerServiceToRepo(serviceMan *serviceModel.Manufacturer) *repoModel.Manufacturer {
	if serviceMan == nil {
		return nil
	}
	return &repoModel.Manufacturer{
		Name:    serviceMan.Name,
		Country: serviceMan.Country,
		Website: serviceMan.Website,
	}
}

func manufacturerRepoToService(repoMan *repoModel.Manufacturer) *serviceModel.Manufacturer {
	if repoMan == nil {
		return nil
	}
	return &serviceModel.Manufacturer{
		Name:    repoMan.Name,
		Country: repoMan.Country,
		Website: repoMan.Website,
	}
}

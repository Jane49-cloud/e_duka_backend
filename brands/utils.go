package brands

import (
	"eleliafrika.com/backend/database"
	"eleliafrika.com/backend/models"
)

func FetchAllBrands() ([]models.Brand, error) {
	var brands []models.Brand
	err := database.Database.Find(&brands).Error
	if err != nil {
		return []models.Brand{}, err
	}

	return brands, nil
}

func FetchSingleBrand(brandname string) (models.Brand, error) {
	var brand models.Brand
	err := database.Database.Where("brand_name=?", brandname).Find(&brand).Error

	if err != nil {
		return models.Brand{}, err
	}
	return brand, nil
}

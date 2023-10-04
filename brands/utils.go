package brands

import (
	"errors"

	"eleliafrika.com/backend/database"
	"eleliafrika.com/backend/models"
)

func FetchAllBrands() ([]models.Brand, error) {
	var brands []models.Brand
	err := database.Database.Where("is_deleted=?", false).Find(&brands).Error
	if err != nil {
		return []models.Brand{}, err
	}

	return brands, nil
}

func FetchSingleBrand(brandname string) (models.Brand, error) {
	var brand models.Brand
	err := database.Database.Where("is_deleted=?", false).Where("brand_name=?", brandname).Find(&brand).Error

	if err != nil {
		return models.Brand{}, err
	}
	return brand, nil
}

func UpdateBrand(brandname string, update models.Brand) (models.Brand, error) {
	var updatedbrand models.Brand
	result := database.Database.Model(updatedbrand).Where("brand_name=?", brandname).Updates(update)
	if result.RowsAffected == 0 {
		return models.Brand{}, errors.New("could not update the brand!! please try again later")
	}
	return updatedbrand, nil

}

package subcategory

import (
	"errors"

	"eleliafrika.com/backend/database"
	"eleliafrika.com/backend/models"
)

func FetchAllSubCategories() ([]models.SubCategory, error) {
	var subcategories []models.SubCategory
	err := database.Database.Where("is_deleted", false).Find(&subcategories).Error
	if err != nil {
		return []models.SubCategory{}, err
	}
	return subcategories, nil
}

func FetchSingleSubCategory(subcategoryname string) (models.SubCategory, error) {
	var subcategory models.SubCategory
	err := database.Database.Where("subcategory_name=?", subcategoryname).Find(&subcategory).Error
	if err != nil {
		return models.SubCategory{}, err
	}
	return subcategory, nil
}

func UpdateSubCategory(subcategoryname string, update models.SubCategory) (models.SubCategory, error) {
	var subcategory models.SubCategory
	result := database.Database.Model(subcategory).Where("subcategory_name=?", subcategoryname).Updates(update)
	if result.RowsAffected == 0 {
		return models.SubCategory{}, errors.New("could not update the category")
	}
	return subcategory, nil
}

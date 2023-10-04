package category

import (
	"errors"

	"eleliafrika.com/backend/database"
	"eleliafrika.com/backend/models"
)

func FetchAllCategories() ([]models.Category, error) {
	var categories []models.Category
	err := database.Database.Where("is_deleted", false).Find(&categories).Error
	if err != nil {
		return []models.Category{}, err
	}
	return categories, nil
}

func FetchSingleCategory(categoryname string) (models.Category, error) {
	var category models.Category
	err := database.Database.Find(&category).Error
	if err != nil {
		return models.Category{}, err
	}
	return category, nil
}

func UpdateCategory(categoryname string, update models.Category) (models.Category, error) {
	var category models.Category
	result := database.Database.Model(category).Where("category_name=?", categoryname).Updates(update)
	if result.RowsAffected == 0 {
		return models.Category{}, errors.New("could not update the brand")
	}
	return category, nil
}

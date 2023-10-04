package category

import (
	"eleliafrika.com/backend/database"
	"eleliafrika.com/backend/models"
)

func FetchAllCategories() ([]models.Category, error) {
	var categories []models.Category
	err := database.Database.Find(&categories).Error
	if err != nil {
		return []models.Category{}, err
	}
	return categories, nil
}

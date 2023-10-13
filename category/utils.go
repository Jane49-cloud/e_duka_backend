package category

import (
	"errors"
	"regexp"

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
	err := database.Database.Where("category_name=?", categoryname).Find(&category).Error
	if err != nil {
		return models.Category{}, err
	}
	return category, nil
}

func UpdateCategory(categoryname string, update models.Category) (models.Category, error) {
	var category models.Category
	result := database.Database.Model(category).Where("category_name=?", categoryname).Updates(update)
	if result.RowsAffected == 0 {
		return models.Category{}, errors.New("could not update the category")
	}
	return category, nil
}

func ValidateCategoryInput(category *models.Category) (bool, error) {
	charPattern := "[!#%^&*()_\\=\\[\\]{};\"\\\\|<>?]"
	if len(category.CategoryName) < 3 {
		return false, errors.New("category name is too short")
	} else if regexp.MustCompile(charPattern).MatchString(category.CategoryName) {
		return false, errors.New("category name cannot contain special characters")
	}
	return true, nil
}

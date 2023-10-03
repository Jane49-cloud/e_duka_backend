package images

import (
	"fmt"

	"eleliafrika.com/backend/database"
	"eleliafrika.com/backend/models"
)

func GetProductImages() ([]models.ProductImage, error) {
	var productsImages []models.ProductImage
	err := database.Database.Find(&productsImages).Error
	if err != nil {
		return []models.ProductImage{}, err

	}
	return productsImages, nil
}

func GetSpecificProductImage(productid string) ([]models.ProductImage, error) {
	var productsImages []models.ProductImage
	fmt.Printf("these are images%s\n", productid)
	err := database.Database.Where("product_id=?", productid).Find(&productsImages).Error
	if err != nil {
		return []models.ProductImage{}, err

	}
	return productsImages, nil
}

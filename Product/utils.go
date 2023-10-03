package product

import (
	"eleliafrika.com/backend/database"
	"eleliafrika.com/backend/models"
)

func FindAllProducts() ([]models.Product, error) {
	var products []models.Product
	err := database.Database.Find(&products).Error
	if err != nil {
		return []models.Product{}, err

	}
	return products, nil
}
func FindSingleProduct(productid string) (models.Product, error) {
	var product models.Product
	err := database.Database.Where("product_id=?", productid).Find(&product).Error
	if err != nil {
		return models.Product{}, err

	}
	return product, nil
}
func Fetchproducts(query string) ([]models.Product, error) {
	var productList []models.Product

	err := database.Database.Where(query).Find(&productList).Error
	if err != nil {
		return []models.Product{}, err
	}
	return productList, nil
}

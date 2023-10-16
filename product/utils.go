package product

import (
	"errors"
	"regexp"
	"strings"
	"unicode"

	"eleliafrika.com/backend/database"
	"eleliafrika.com/backend/images"
)

func FindSingleProduct(query string) (Product, error) {
	var product Product
	err := database.Database.Where("product_id=?", query).Find(&product).Error
	if err != nil {
		return Product{}, err

	}
	return product, nil
}
func Fetchproducts() ([]Product, error) {
	var productList []Product

	err := database.Database.Find(&productList).Error
	if err != nil {
		return []Product{}, err
	}
	if len(productList) > 0 {
		for i, product := range productList {
			mainImage, err := images.DownloadImageFromBucket(product.MainImage)
			if err != nil {
				return productList, err
			} else if product.MainImage == "" {
				return productList, errors.New("could not download image from the storage")
			}

			productList[i].MainImage = mainImage
		}
	}
	return productList, nil
}

func FetchAds() ([]Product, error) {
	var productList []Product

	err := database.Database.Where("is_deleted=?", false).Where("is_approved=?", true).Where("is_active=?", true).Where("is_suspended=?", false).Find(&productList).Error
	if err != nil {
		return []Product{}, err
	}
	if len(productList) > 0 {
		for i, product := range productList {
			mainImage, err := images.DownloadImageFromBucket(product.MainImage)
			if err != nil {
				return productList, err
			} else if product.MainImage == "" {
				return productList, errors.New("could not download image from the storage")
			}

			productList[i].MainImage = mainImage
		}
	}
	return productList, nil
}

func ValidateProductInput(product *AddProductInput) (bool, error) {
	productDetails := []string{product.ProductName, product.ProductPrice, product.ProductDescription, product.MainImage, product.ProductType, product.Brand, product.Category, product.SubCategory}
	charPattern := "[!@#$%^&*()\\=\\[\\]{};\\\\|<>?]"
	for _, value := range productDetails {
		if value == product.ProductName {
			value = strings.TrimSpace(value)
			if len(value) < 3 {
				return false, errors.New("product name should be atleast three characters long")
			} else if regexp.MustCompile(charPattern).MatchString(product.ProductName) {
				return false, errors.New("product name should not contain special character")
			}
		} else if value == product.ProductPrice {
			if len(value) < 2 {
				return false, errors.New("price invalid")
			} else {
				for _, char := range product.ProductPrice {
					if !unicode.IsNumber(char) {
						return false, errors.New("price can only contain numbers")
					}
				}
			}
		} else if value == product.ProductDescription {
			charPattern := "[@#$%^&\\=\\[\\]{};:\\\\|<>]"
			value = strings.TrimSpace(value)
			if len(value) < 100 {
				return false, errors.New("product description should atleast be 100 characters long")
			} else if regexp.MustCompile(charPattern).MatchString(product.ProductDescription) {
				return false, errors.New("product description should not contain special character")
			}
		} else if value == product.MainImage {
			value = strings.TrimSpace(value)
			if value == "" {
				return false, errors.New("product image cannot be empty")
			}

		} else if value == product.ProductType {
			value = strings.TrimSpace(value)
			if len(value) < 3 {
				return false, errors.New("product type should atleast be 3 characters long")
			} else if regexp.MustCompile(charPattern).MatchString(product.ProductType) {
				return false, errors.New("product type should not contain special character")
			}
		} else if value == product.Brand {
			value = strings.TrimSpace(value)
			if value == "" {
				return false, errors.New("product brand should not be empty")
			} else if regexp.MustCompile(charPattern).MatchString(product.Brand) {
				return false, errors.New("product brand should not contain special character")
			}
		} else if value == product.Category {
			value = strings.TrimSpace(value)
			if len(value) < 3 {
				return false, errors.New("product category should atleast be 3 characters long")
			} else if regexp.MustCompile(charPattern).MatchString(product.Category) {
				return false, errors.New("product category should not contain special character")
			}
		} else if value == product.SubCategory {
			value = strings.TrimSpace(value)
			if len(value) < 3 {
				return false, errors.New("product sub category should atleast be 3 characters long")
			} else if regexp.MustCompile(charPattern).MatchString(product.SubCategory) {
				return false, errors.New("product sub category should not contain special character")
			}
		}
	}
	return true, nil
}
func UpdateProductUtil(query string, update Product) (Product, error) {
	var updatedProduct Product

	result := database.Database.Model(&updatedProduct).Where(query).Updates(update)

	if result.RowsAffected == 0 {
		return Product{}, errors.New("could not update the product right now")
	}
	return updatedProduct, nil
}

func ActivateProductUtil(query string) (bool, error) {
	var updatedProduct Product
	result := database.Database.Model(&updatedProduct).Where(query).Update("is_active", true)
	if result.RowsAffected == 0 {
		return false, errors.New("could not activate the current product")
	}
	return true, nil
}
func DeactivateProductUtil(query string) (bool, error) {
	var updatedProduct Product
	result := database.Database.Model(&updatedProduct).Where(query).Update("is_active", false)
	if result.RowsAffected == 0 {
		return false, errors.New("could not deactivate the current product")
	}
	return true, nil
}

func DeleteProductUtil(query string) (bool, error) {
	var updatedProduct Product
	result := database.Database.Model(&updatedProduct).Where(query).Update("is_deleted", true)
	if result.RowsAffected == 0 {
		return false, errors.New("could not delete the current product")
	}
	return true, nil
}
func RestoreProductUtil(query string) (bool, error) {
	var updatedProduct Product
	result := database.Database.Model(&updatedProduct).Where(query).Update("is_deleted", false)
	if result.RowsAffected == 0 {
		return false, errors.New("could not restore the current product")
	}
	return true, nil
}

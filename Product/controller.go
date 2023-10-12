package product

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"eleliafrika.com/backend/category"
	globalutils "eleliafrika.com/backend/global_utils"
	"eleliafrika.com/backend/images"
	"eleliafrika.com/backend/models"
	subcategory "eleliafrika.com/backend/subcategories"
	"eleliafrika.com/backend/users"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AddProduct(context *gin.Context) {
	var productInput AddProductInput

	if err := context.ShouldBind(&productInput); err != nil {
		response := models.Reply{
			Message: "could not bind data from the user",
			Error:   err.Error(),
			Success: false,
		}
		context.JSON(http.StatusBadRequest, response)
		return
	}

	productuuid := uuid.New()
	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02 15:04:05")
	success, err := ValidateProductInput(&productInput)
	if err != nil {
		response := models.Reply{
			Message: "error validating user input",
			Error:   err.Error(),
			Success: false,
			Data:    productInput,
		}
		context.JSON(http.StatusBadRequest, response)
		return
	} else if !success {
		response := models.Reply{
			Message: "error validating user input for sub category",
			Success: false,
		}
		context.JSON(http.StatusBadRequest, response)
		return
	} else {

		// get current user to add to userid field
		user, err := users.CurrentUser(context)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "could not find user",
			})
			return
		} else if user.Firstname == "" {
			globalutils.HandleSuccess("user not found", user, context)
		} else {

			// check if category exists
			categoryExists, err := category.FetchSingleCategory(productInput.Category)
			if err != nil {
				response := models.Reply{
					Message: "error validating the category",
					Error:   err.Error(),
					Success: false,
				}
				context.JSON(http.StatusBadRequest, response)
				return
			}
			subCategoryExists, err := subcategory.FetchSingleSubCategory(productInput.SubCategory)
			if err != nil {
				response := models.Reply{
					Message: "error validating the sub category",
					Error:   err.Error(),
					Success: false,
				}
				context.JSON(http.StatusBadRequest, response)
				return
			}

			if categoryExists.CategoryName == "" {
				response := models.Reply{
					Message: "category not found",
					Success: false,
				}
				context.JSON(http.StatusBadRequest, response)
				return
			} else if subCategoryExists.SubCategoryName == "" {
				response := models.Reply{
					Message: "sub category not found",
					Success: false,
				}
				context.JSON(http.StatusBadRequest, response)
				return
			} else {

				// handle image input
				mainImagePath, err := images.UploadMainimage(productInput.MainImage, productInput.ProductName)

				if err != nil {
					globalutils.HandleError("error uploading main image", err, context)
					return
				}
				product := models.Product{
					ProductID:          productuuid.String(),
					ProductName:        productInput.ProductName,
					ProductPrice:       productInput.ProductPrice,
					ProductDescription: productInput.ProductDescription,
					UserID:             user.UserID,
					MainImage:          mainImagePath,
					ProductStatus:      "Active",
					Quantity:           productInput.Quantity,
					ProductType:        productInput.ProductType,
					TotalLikes:         0,
					TotalComments:      0,
					DateAdded:          formattedTime,
					LastUpdated:        formattedTime,
					LatestInteractions: formattedTime,
					TotalInteractions:  0,
					TotalBookmarks:     0,
					Brand:              productInput.Brand,
					Category:           productInput.Category,
					SubCategory:        productInput.SubCategory,
				}
				imagesPath, err := images.UploadOtherImages(productInput.ProductImages, product.ProductName)
				if err != nil {
					globalutils.HandleError("error uploading product images", err, context)
					return
				} else {
					for _, i := range imagesPath {

						imageuuid := uuid.New()
						image := models.ProductImage{
							ImageID:   imageuuid.String(),
							ProductID: productuuid.String(),
							ImageUrl:  i,
						}
						savedImage, err := image.Save()
						if err != nil {
							context.JSON(http.StatusBadRequest, gin.H{
								"error with saving image": err.Error(),
								"success":                 false,
								"image":                   savedImage,
							})
							return
						}
					}
				}

				savedProduct, err := product.Save()

				if err != nil {
					context.JSON(http.StatusBadRequest, gin.H{
						"error with save": err.Error(),
						"success":         false,
					})
					return
				}
				imageString, err := images.DownloadImageFromBucket(mainImagePath)

				if err != nil {
					context.JSON(http.StatusBadRequest, gin.H{
						"error with downloading image": err.Error(),
						"success":                      false,
					})
					return
				}
				savedProduct.MainImage = imageString
				context.JSON(http.StatusCreated, gin.H{
					"data":    savedProduct,
					"success": true,
					"message": "product has been added succesfully",
				})
			}
		}
	}
}
func GetAllProducts(context *gin.Context) {
	var err error

	products, err := Fetchproducts()
	if err != nil {
		response := models.Reply{
			Message: "Error fetching product",
			Error:   err.Error(),
			Success: false,
		}
		context.JSON(http.StatusBadRequest, response)
	} else {
		response := models.Reply{
			Message: "all products fetched",
			Success: true,
			Data:    products,
		}
		context.JSON(http.StatusOK, response)
		return
	}

}
func GetSingleProduct(context *gin.Context) {

	productid := context.Query("id")

	query := "product_id=" + productid

	productExist, err := FindSingleProduct(query)
	if err != nil {
		globalutils.HandleError("could not fetch single product", err, context)
		return
	} else if productExist.ProductName != "" {
		newId := strings.ReplaceAll(productid, "'", "")
		productImages, err := images.GetSpecificProductImage(newId)
		if err != nil {
			globalutils.HandleError("could not fetch product images", err, context)
		}

		productData := gin.H{
			"productdata": productExist,
			"images":      productImages,
		}
		globalutils.HandleSuccess("fetched product succesful", productData, context)
		return
	} else {
		globalutils.HandleSuccess("fetch product does not exist", productExist, context)
		return
	}
}
func UpdateProduct(context *gin.Context) {

	var productUpdate AddProductInput
	if err := context.ShouldBindJSON(&productUpdate); err != nil {
		response := models.Reply{
			Message: "could not bind the user data to the request needs",
			Error:   err.Error(),
			Success: false,
			Data:    productUpdate,
		}
		context.JSON(http.StatusBadRequest, response)
		return
	}

	success, err := ValidateProductInput(&productUpdate)

	if err != nil {
		response := models.Reply{
			Message: "error validating user input",
			Error:   err.Error(),
			Success: false,
		}
		context.JSON(http.StatusBadRequest, response)
		return
	} else if !success {
		response := models.Reply{
			Message: "error validating user input for product data",
			Success: false,
		}
		context.JSON(http.StatusBadRequest, response)
		return
	} else {
		user, err := users.CurrentUser(context)
		if err != nil {
			response := models.Reply{
				Message: "could not get current user!!user required in order to update product",
				Error:   err.Error(),
				Success: false,
			}
			context.JSON(http.StatusBadRequest, response)
			return
		} else if user.Firstname == "" {
			response := models.Reply{
				Message: "no such user",
				Success: false,
			}
			context.JSON(http.StatusBadRequest, response)
			return
		} else {

			productid := context.Query("id")
			if productid != "" {
				query := "product_id=" + productid
				productExist, err := FindSingleProduct(query)
				if err != nil {
					response := models.Reply{
						Message: "could not fetch the product",
						Error:   err.Error(),
						Success: false,
					}
					context.JSON(http.StatusBadRequest, response)
					return
				} else if productExist.ProductName == "" {
					response := models.Reply{
						Message: "product does not exist",
						Success: false,
					}
					context.JSON(http.StatusBadRequest, response)
					return
				} else {
					newproduct := models.Product{
						ProductName:        productUpdate.ProductName,
						ProductPrice:       productUpdate.ProductPrice,
						ProductDescription: productUpdate.ProductDescription,
						MainImage:          productUpdate.MainImage,
						Quantity:           productUpdate.Quantity,
						ProductType:        productUpdate.ProductType,
						Brand:              productUpdate.Brand,
						Category:           productUpdate.Category,
						SubCategory:        productUpdate.SubCategory,
					}
					productUpdated, err := UpdateProductUtil(query, newproduct)

					if err != nil {
						response := models.Reply{
							Message: "could not update product",
							Error:   err.Error(),
							Success: false,
						}
						context.JSON(http.StatusBadRequest, response)
						return
					} else if productUpdated.ProductName == "" {
						response := models.Reply{
							Message: "could not update product",
							Success: false,
						}
						context.JSON(http.StatusBadRequest, response)
						return
					} else {
						response := models.Reply{
							Message: "Product updated",
							Success: true,
							Data:    productUpdated,
						}
						context.JSON(http.StatusOK, response)
						return
					}
				}
			} else {
				response := models.Reply{
					Message: "Cannot send request without header!!Invalid request",
					Success: false,
					Error:   "supply product id!!",
				}
				context.JSON(http.StatusBadRequest, response)
			}
		}
	}
}

func ActivateProduct(context *gin.Context) {
	productid := context.Query("id")
	query := "product_id=" + productid

	// check if product exist
	productExist, err := FindSingleProduct(query)
	if err != nil {
		globalutils.HandleError("error finding product", err, context)
	} else if productExist.ProductName == "" {
		globalutils.HandleSuccess("the product does not exist", models.Product{}, context)
	} else if productExist.IsDeleted {
		globalutils.HandleSuccess("cannot activate a deleted product!!Please restore product first", productExist, context)
	} else if productExist.IsActive {
		globalutils.HandleSuccess("product is already active", productExist, context)
	} else {
		success, err := ActivateProductUtil(query)
		if err != nil {
			globalutils.HandleError("error activating  product", err, context)
		} else if !success {
			globalutils.HandleError("failed in activating product", errors.New("could not activate product!!try again"), context)
		} else {
			globalutils.HandleSuccess("succesfully activated the product", productExist, context)
		}
	}

}

func DeactivateProduct(context *gin.Context) {
	productid := context.Query("id")
	query := "product_id=" + productid

	// check if product exist
	productExist, err := FindSingleProduct(query)
	if err != nil {
		globalutils.HandleError("error finding product", err, context)
	} else if productExist.ProductName == "" {
		globalutils.HandleSuccess("the product does not exist", models.Product{}, context)
	} else if productExist.IsDeleted {
		globalutils.HandleSuccess("product is deleted!!Please restore product first", productExist, context)
	} else if !productExist.IsActive {
		globalutils.HandleSuccess("product is not active active", productExist, context)
	} else {
		success, err := DeactivateProductUtil(query)
		if err != nil {
			globalutils.HandleError("error deactivating  product", err, context)
		} else if !success {
			globalutils.HandleError("failed in deactivating product", errors.New("could not deactivate product!!try again"), context)
		} else {
			globalutils.HandleSuccess("succesfully deactivated the product", productExist, context)
		}
	}

}

func DeleteProduct(context *gin.Context) {
	productid := context.Query("id")
	query := "product_id=" + productid

	// check if product exist
	productExist, err := FindSingleProduct(query)
	if err != nil {
		globalutils.HandleError("error finding product", err, context)
	} else if productExist.ProductName == "" {
		globalutils.HandleSuccess("the product does not exist", models.Product{}, context)
	} else if productExist.IsDeleted {
		globalutils.HandleSuccess("product already deleted", productExist, context)
	} else if productExist.IsActive {
		globalutils.HandleSuccess("you cannot delete an active product!! Please deactivate the product first", productExist, context)
	} else {
		success, err := DeleteProductUtil(query)
		if err != nil {
			globalutils.HandleError("error deleting  product", err, context)
		} else if !success {
			globalutils.HandleError("failed in deleting product", errors.New("could not delete product!!try again"), context)
		} else {
			globalutils.HandleSuccess("succesfully deleted the product", productExist, context)
		}
	}

}

func RestoreProduct(context *gin.Context) {
	productid := context.Query("id")
	query := "product_id=" + productid

	// check if product exist
	productExist, err := FindSingleProduct(query)
	if err != nil {
		globalutils.HandleError("error finding product", err, context)
	} else if productExist.ProductName == "" {
		globalutils.HandleSuccess("the product does not exist", models.Product{}, context)
	} else if !productExist.IsDeleted {
		globalutils.HandleSuccess("product is not deleted", productExist, context)
	} else {
		success, err := RestoreProductUtil(query)
		if err != nil {
			globalutils.HandleError("error restoring  product", err, context)
		} else if !success {
			globalutils.HandleError("failed in restoring product", errors.New("could not restore product!!try again"), context)
		} else {
			globalutils.HandleSuccess("succesfully restoring the product", productExist, context)
		}
	}

}

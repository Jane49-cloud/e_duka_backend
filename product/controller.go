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
			Success: false,
			Error:   err.Error(),
		}
		context.JSON(http.StatusOK, response)
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
			response := models.Reply{
				Message: "error fetching user",
				Success: false,
				Error:   err.Error(),
			}
			context.JSON(http.StatusBadRequest, response)
			return
		} else if user.Firstname == "" {
			response := models.Reply{
				Message: "user not found",
				Success: false,
				Error:   err.Error(),
			}
			context.JSON(http.StatusBadRequest, response)
			return

			// else if !user.IsApproved {
			// 	globalutils.UnAuthorized(context)
			// }
		} else {

			// check if category exists
			categoryExists, err := category.FetchSingleCategory(productInput.Category)
			if err != nil {
				response := models.Reply{
					Message: "error validating category",
					Success: false,
					Error:   errors.New("error validating category").Error(),
				}
				context.JSON(http.StatusBadRequest, response)
				return
			}
			subCategoryExists, err := subcategory.FetchSingleSubCategory(productInput.SubCategory)
			if err != nil {
				response := models.Reply{
					Message: "error validating sub category",
					Success: false,
					Error:   errors.New("error validating sub category").Error(),
				}
				context.JSON(http.StatusBadRequest, response)
				return
			}

			if categoryExists.CategoryName == "" {
				response := models.Reply{
					Message: "category not found",
					Success: false,
					Error:   errors.New("category not found").Error(),
				}
				context.JSON(http.StatusBadRequest, response)
				return
			} else if subCategoryExists.SubCategoryName == "" {
				response := models.Reply{
					Message: "sub category not found",
					Success: false,
					Error:   errors.New("sub category not found").Error(),
				}
				context.JSON(http.StatusBadRequest, response)
				return
			} else {
				product := Product{
					ProductID:          productuuid.String(),
					ProductName:        productInput.ProductName,
					ProductPrice:       productInput.ProductPrice,
					ProductDescription: productInput.ProductDescription,
					UserID:             user.UserID,
					MainImage:          productInput.MainImage,
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

				savedProduct, err := product.Save()

				if err != nil {
					response := models.Reply{
						Message: "error uploading product images",
						Success: false,
						Error:   err.Error(),
					}
					context.JSON(http.StatusBadRequest, response)
					return
				} else {
					for _, i := range productInput.ProductImages {

						imageuuid := uuid.New()
						image := models.ProductImage{
							ImageID:   imageuuid.String(),
							ProductID: productuuid.String(),
							ImageUrl:  i,
						}
						_, err := image.Save()
						if err != nil {
							response := models.Reply{
								Message: "error with saving image",
								Success: false,
								Error:   err.Error(),
							}
							context.JSON(http.StatusBadRequest, response)
							return
						}
					}
				}

				response := models.Reply{
					Message: "product has been added succesfully",
					Success: true,
					Data:    savedProduct,
				}
				context.JSON(http.StatusCreated, response)
				return
			}
		}
	}
}
func GetAllProducts(context *gin.Context) {

	products, err := Fetchproducts()
	if err != nil {
		response := models.Reply{
			Message: "error fetching products",
			Success: false,
			Error:   err.Error(),
		}
		context.JSON(http.StatusBadRequest, response)
		return
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
func GetAllAds(context *gin.Context) {
	var err error

	products, err := FetchAds()
	if err != nil {
		response := models.Reply{
			Message: "error fetching ads",
			Success: false,
			Error:   err.Error(),
		}
		context.JSON(http.StatusBadRequest, response)
		return
	} else {
		response := models.Reply{
			Message: "all ads fetched",
			Success: true,
			Data:    products,
		}
		context.JSON(http.StatusOK, response)
		return
	}

}
func GetSingleProduct(context *gin.Context) {

	productid := context.Param("id")

	productExist, err := FindSingleProduct(productid)
	if err != nil {
		response := models.Reply{
			Message: "could not fetch single product",
			Success: false,
			Error:   err.Error(),
		}
		context.JSON(http.StatusBadRequest, response)
		return
	} else if productExist.ProductName != "" {
		newId := strings.ReplaceAll(productid, "'", "")
		productImages, err := images.GetSpecificProductImage(newId)
		if err != nil {
			globalutils.HandleError("could not download product images", err, context)
		}
		mainImage, err := images.DownloadImageFromBucket(productExist.MainImage)
		if err != nil {
			globalutils.HandleError("could not download product main image", err, context)
		}
		productExist.MainImage = mainImage

		// fetch user details of the product owner
		currentuser, err := users.FindUserById(string(productExist.UserID))
		if err != nil {
			globalutils.HandleError("error finding the seller", err, context)
			return
		}
		sellerDetails := gin.H{
			"seller_name":        currentuser.Firstname + " " + currentuser.Lastname,
			"seller_email":       currentuser.Email,
			"seller_phonenumber": currentuser.Phone,
			"seller_location":    currentuser.Location,
			"user_profile":       currentuser.UserImage,
		}

		productData := gin.H{
			"productdata":    productExist,
			"images":         productImages,
			"seller_details": sellerDetails,
		}
		globalutils.HandleSuccess("fetched product succesful", productData, context)
		return
	} else {
		globalutils.HandleSuccess("fetch product does not exist", productExist, context)
		return
	}
}
func GetSingleAd(context *gin.Context) {

	productid := context.Param("id")

	// query := "product_id=" + productid
	productExist, err := FindSingleAd(productid)
	if err != nil {
		globalutils.HandleError("could not fetch single ad", err, context)
		return
	} else if productExist.ProductName != "" {
		newId := strings.ReplaceAll(productid, "'", "")
		productImages, err := images.GetSpecificProductImage(newId)
		if err != nil {
			globalutils.HandleError("could not download ad images", err, context)
		}
		mainImage, err := images.DownloadImageFromBucket(productExist.MainImage)
		if err != nil {
			globalutils.HandleError("could not download ad main image", err, context)
		}
		productExist.MainImage = mainImage

		// fetch user details of the product owner
		currentuser, err := users.FindUserById(string(productExist.UserID))
		if err != nil {
			globalutils.HandleError("error finding the seller", err, context)
			return
		}
		sellerDetails := gin.H{
			"seller_name":        currentuser.Firstname + " " + currentuser.Lastname,
			"seller_email":       currentuser.Email,
			"seller_phonenumber": currentuser.Phone,
			"seller_location":    currentuser.Location,
			"user_profile":       currentuser.UserImage,
		}

		productData := gin.H{
			"productdata":    productExist,
			"images":         productImages,
			"seller_details": sellerDetails,
		}
		globalutils.HandleSuccess("fetched ad succesful", productData, context)
		return
	} else {
		globalutils.HandleSuccess("fetch ad does not exist", productExist, context)
		return
	}
}
func UpdateProduct(context *gin.Context) {

	var productUpdate AddProductInput
	if err := context.ShouldBindJSON(&productUpdate); err != nil {
		globalutils.HandleError("could not bind the user data to the request needs", err, context)
		return
	}

	success, err := ValidateProductInput(&productUpdate)

	if err != nil {

		globalutils.HandleError("error validating user input", err, context)
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

			globalutils.HandleError("could not get current user!!user required in order to update product", err, context)
			return
		} else if user.Firstname == "" {
			globalutils.UnAuthenticated(context)
			return
		} else {

			productid := context.Query("id")
			if productid != "" {

				query := "product_id=" + productid
				id := strings.ReplaceAll(productid, "'", "")

				productExist, err := FindSingleProduct(id)
				if err != nil {
					globalutils.HandleError("could not fetch product", err, context)
					return
				}
				userOwnsProduct, err := ValidateUserOwnsProduct(user.UserID, productExist.UserID)
				if err != nil {
					globalutils.HandleError("error occurred while validating user", err, context)
				} else if !userOwnsProduct {
					globalutils.UnAuthorized(context)
				} else if productExist.ProductName == "" {

					globalutils.HandleError("product does not exist", errors.New("error fetching the product"), context)
					return
				} else {
					newproduct := Product{
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

	id := strings.ReplaceAll(productid, "'", "")

	// check if product exist
	productExist, err := FindSingleProduct(id)
	if err != nil {
		globalutils.HandleError("error finding product", err, context)
	} else if productExist.ProductName == "" {
		globalutils.HandleSuccess("the product does not exist", Product{}, context)
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

	id := strings.ReplaceAll(productid, "'", "")

	// check if product exist
	productExist, err := FindSingleProduct(id)
	if err != nil {
		globalutils.HandleError("error finding product", err, context)
	} else if productExist.ProductName == "" {
		globalutils.HandleSuccess("the product does not exist", Product{}, context)
	} else if productExist.IsDeleted {
		globalutils.HandleSuccess("product is deleted!!Please restore product first", productExist, context)
	} else if !productExist.IsActive {
		globalutils.HandleSuccess("product is not active", productExist, context)
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

	id := strings.ReplaceAll(productid, "'", "")

	// check if product exist
	productExist, err := FindSingleProduct(id)
	if err != nil {
		globalutils.HandleError("error finding product", err, context)
	} else if productExist.ProductName == "" {
		globalutils.HandleSuccess("the product does not exist", Product{}, context)
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

	id := strings.ReplaceAll(productid, "'", "")

	// check if product exist
	productExist, err := FindSingleProduct(id)
	if err != nil {
		globalutils.HandleError("error finding product", err, context)
	} else if productExist.ProductName == "" {
		globalutils.HandleSuccess("the product does not exist", Product{}, context)
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

func FetchSingleUserProducts(context *gin.Context) {
	id := context.Query("id")
	products, err := FetchSingleUserProductsUtil(strings.ReplaceAll(id, "'", ""))
	if err != nil {
		response := models.Reply{
			Message: "error fetching single user products",
			Success: false,
			Error:   err.Error(),
		}
		context.JSON(http.StatusBadRequest, response)
		return
	} else {
		response := models.Reply{
			Message: "single user products fetched",
			Success: true,
			Data:    products,
		}
		context.JSON(http.StatusOK, response)
		return
	}
}

func FetchSingleUserAds(context *gin.Context) {
	id := context.Query("id")

	products, err := FetchSingleUserAdsUtil(strings.ReplaceAll(id, "'", ""))
	if err != nil {
		globalutils.HandleError("error fetching single user products", err, context)
		return
	} else {
		globalutils.HandleSuccess("single user products fetched", products, context)
		return
	}
}

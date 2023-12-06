package product

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"eleliafrika.com/backend/category"
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
			Message: err.Error(),
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
			Error:   errors.New("operation not succesfull").Error(),
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
			context.JSON(http.StatusUnauthorized, response)
			return
		} else if user.Firstname == "" {
			response := models.Reply{
				Message: "user not found",
				Success: false,
				Error:   errors.New("user not found").Error(),
			}
			context.JSON(http.StatusUnauthorized, response)
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

				imageUrl, err := images.UploadHandler(productInput.ProductName, productInput.MainImage, context)
				if err != nil {
					response := models.Reply{
						Message: "main image not saved",
						Success: false,
						Error:   err.Error(),
					}
					context.JSON(http.StatusBadRequest, response)
					return
				}
				product := Product{
					ProductID:          productuuid.String(),
					ProductName:        productInput.ProductName,
					ProductPrice:       productInput.ProductPrice,
					ProductDescription: productInput.ProductDescription,
					UserID:             user.UserID,
					MainImage:          imageUrl,
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
						imageUrl, err := images.UploadHandler(productInput.ProductName, i, context)
						if err != nil {
							response := models.Reply{
								Message: "error with saving image",
								Success: false,
								Error:   err.Error(),
							}
							context.JSON(http.StatusBadRequest, response)
							return
						}
						imageuuid := uuid.New()
						image := models.ProductImage{
							ImageID:   imageuuid.String(),
							ProductID: productuuid.String(),
							ImageUrl:  imageUrl,
						}
						_, err = image.Save()
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
	query := context.Query("search")
	query = strings.ReplaceAll(strings.ToLower(query), "'", "")

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
		var productList []Product
		addedProducts := make(map[uint]bool)

		if query != "" {
			for _, item := range products {
				splitName := strings.ReplaceAll(query, "and", "")
				fmt.Println(splitName)
				newQ := strings.Split(splitName, " ")
				for _, segment := range newQ {
					if segment != "" {
						if strings.ToLower(item.ProductName) == segment || strings.Contains(strings.ToLower(item.ProductName), segment) {
							if _, exists := addedProducts[item.ID]; !exists {
								productList = append(productList, item)
								addedProducts[item.ID] = true
							}
						} else if strings.ToLower(item.Category) == segment || strings.Contains(strings.ToLower(item.Category), segment) {
							if _, exists := addedProducts[item.ID]; !exists {
								productList = append(productList, item)
								addedProducts[item.ID] = true
							}
						} else if strings.ToLower(item.SubCategory) == segment || strings.Contains(strings.ToLower(item.SubCategory), segment) {
							if _, exists := addedProducts[item.ID]; !exists {
								productList = append(productList, item)
								addedProducts[item.ID] = true
							}
						} else if strings.ToLower(item.Brand) == segment || strings.Contains(strings.ToLower(item.Brand), segment) {
							if _, exists := addedProducts[item.ID]; !exists {
								productList = append(productList, item)
								addedProducts[item.ID] = true
							}
						} else if strings.Contains(strings.ToLower(item.ProductDescription), segment) {
							if _, exists := addedProducts[item.ID]; !exists {
								productList = append(productList, item)
								addedProducts[item.ID] = true
							}
						}
					}

				}

				currentuser, err := users.FindUserById(string(item.UserID))
				if err != nil {
					response := models.Reply{
						Error:   err.Error(),
						Message: "error finding the seller",
						Success: false,
					}
					context.JSON(http.StatusBadRequest, response)
					return
				} else {
					fname := currentuser.Firstname
					mName := currentuser.Middlename
					lname := currentuser.Lastname

					if strings.ToLower(fname) == query || strings.Contains(fname, query) {
						productList = append(productList, item)
					} else if strings.ToLower(lname) == query || strings.Contains(lname, query) {
						productList = append(productList, item)
					} else if strings.ToLower(mName) == query || strings.Contains(mName, query) {
						productList = append(productList, item)
					}
				}

			}
		} else {
			productList = products
		}

		var data []interface{}
		for _, product := range productList {

			currentuser, err := users.FindUserById(string(product.UserID))
			if err != nil {
				response := models.Reply{
					Error:   err.Error(),
					Message: "error finding the seller",
					Success: false,
				}
				context.JSON(http.StatusBadRequest, response)
				return
			} else {
				productData := gin.H{
					"product_data": product,
					"user_name":    currentuser.Firstname + " " + currentuser.Middlename + " " + currentuser.Lastname,
				}
				data = append(data, productData)

			}

		}
		response := models.Reply{
			Message: "all ads fetched",
			Success: true,
			Data:    data,
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
			response := models.Reply{
				Error:   err.Error(),
				Message: "could not download product images",
				Success: false,
			}
			context.JSON(http.StatusBadRequest, response)
			return

		}
		// fetch user details of the product owner
		currentuser, err := users.FindUserById(string(productExist.UserID))
		if err != nil {
			response := models.Reply{
				Error:   err.Error(),
				Message: "error finding the seller",
				Success: false,
			}
			context.JSON(http.StatusBadRequest, response)
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
		response := models.Reply{
			Data:    productData,
			Message: "fetched product succesful",
			Success: true,
		}
		context.JSON(http.StatusOK, response)
		return
	} else {
		response := models.Reply{
			Error:   err.Error(),
			Message: "fetch product does not exist",
			Success: true,
		}
		context.JSON(http.StatusBadRequest, response)
		return
	}
}
func GetSingleAd(context *gin.Context) {

	productid := context.Param("id")

	// query := "product_id=" + productid
	productExist, err := FindSingleAd(productid)
	if err != nil {
		response := models.Reply{
			Message: "could not fetch single ad",
			Success: false,
			Error:   err.Error(),
		}
		context.JSON(http.StatusBadRequest, response)
		return
	} else if productExist.ProductName != "" {

		// fetch user details of the product owner
		currentuser, err := users.FindUserById(string(productExist.UserID))
		if err != nil {
			response := models.Reply{
				Message: "error finding the seller",
				Success: false,
				Error:   err.Error(),
			}
			context.JSON(http.StatusBadRequest, response)
			return
		}

		images, err := images.GetSpecificProductImage(productid)
		if err != nil {
			response := models.Reply{
				Message: "could not get the images",
				Success: false,
				Error:   err.Error(),
			}
			context.JSON(http.StatusBadRequest, response)
			return
		}
		similarProducts, err := FetchSimilarProducts(productExist.Category)
		if err != nil {
			response := models.Reply{
				Error:   err.Error(),
				Message: "error finding the data",
				Success: false,
			}
			context.JSON(http.StatusBadRequest, response)
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
			"product_data":     productExist,
			"seller_details":   sellerDetails,
			"product_images":   images,
			"similar_products": similarProducts,
		}

		if err != nil {
			response := models.Reply{
				Message: "could not get the images",
				Success: false,
				Error:   err.Error(),
			}
			context.JSON(http.StatusBadRequest, response)
			return
		}
		response := models.Reply{
			Message: "fetched single ad",
			Success: true,
			Data:    productData,
		}
		context.JSON(http.StatusOK, response)
		return
	} else {
		response := models.Reply{
			Message: "fetched ad does not exist",
			Success: false,
			Error:   errors.New("ad does not exist").Error(),
		}
		context.JSON(http.StatusBadRequest, response)
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
				Error:   err.Error(),
				Message: "error finding user",
				Success: false,
			}
			context.JSON(http.StatusBadRequest, response)
			return
		} else {

			productid := context.Query("id")
			if productid != "" {

				query := "product_id=" + productid
				id := strings.ReplaceAll(productid, "'", "")

				productExist, err := FindSingleProduct(id)
				if err != nil {
					response := models.Reply{
						Error:   err.Error(),
						Message: "could not fetch product",
						Success: false,
					}
					context.JSON(http.StatusBadRequest, response)
					return
				}
				userOwnsProduct, err := ValidateUserOwnsProduct(user.UserID, productExist.UserID)
				if err != nil {
					response := models.Reply{
						Error:   err.Error(),
						Message: "error occurred while validating user",
						Success: false,
					}
					context.JSON(http.StatusBadRequest, response)
					return
				} else if !userOwnsProduct {
					response := models.Reply{
						Error:   err.Error(),
						Message: "error authorizing operation user",
						Success: false,
					}
					context.JSON(http.StatusBadRequest, response)
					return
				} else if productExist.ProductName == "" {
					response := models.Reply{
						Error:   errors.New("error fetching the product").Error(),
						Message: "product does not exist",
						Success: false,
					}
					context.JSON(http.StatusBadRequest, response)
					return
				} else {
					var imageUrl = ""
					if productUpdate.MainImage != "" {
						imageUrl, err = images.UploadHandler(productUpdate.ProductName, productUpdate.MainImage, context)
						if err != nil {
							response := models.Reply{
								Message: "could not update product",
								Error:   err.Error(),
								Success: false,
							}
							context.JSON(http.StatusBadRequest, response)
							return
						}
					} else {
						imageUrl = productExist.MainImage
					}

					newproduct := Product{
						ProductName:        productUpdate.ProductName,
						ProductPrice:       productUpdate.ProductPrice,
						ProductDescription: productUpdate.ProductDescription,
						MainImage:          imageUrl,
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
							Error:   errors.New("could not update product").Error(),
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
				context.JSON(http.StatusOK, response)
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
		response := models.Reply{
			Error:   err.Error(),
			Message: "error finding product",
			Success: false,
		}
		context.JSON(http.StatusBadRequest, response)
		return

	} else if productExist.ProductName == "" {
		response := models.Reply{
			Error:   errors.New("product does not exist").Error(),
			Message: "the product does not exist",
			Success: false,
		}
		context.JSON(http.StatusOK, response)
		return

	} else if productExist.IsDeleted {
		response := models.Reply{
			Error:   errors.New("product is deleted").Error(),
			Message: "cannot activate a deleted product!!Please restore product first",
			Success: false,
		}
		context.JSON(http.StatusBadRequest, response)
		return

	} else if productExist.IsActive {
		response := models.Reply{
			Error:   err.Error(),
			Message: "product is already active",
			Success: false,
		}
		context.JSON(http.StatusBadRequest, response)
		return

	} else {
		success, err := ActivateProductUtil(query)
		if err != nil {
			response := models.Reply{
				Error:   err.Error(),
				Message: "product is already active",
				Success: false,
			}
			context.JSON(http.StatusBadRequest, response)
			return
		} else if !success {
			response := models.Reply{
				Error:   errors.New("could not activate product").Error(),
				Message: "failed in activating product",
				Success: false,
			}
			context.JSON(http.StatusBadRequest, response)
			return
		} else {
			response := models.Reply{
				Data:    productExist,
				Message: "succesfully activated the product",
				Success: true,
			}
			context.JSON(http.StatusOK, response)
			return

		}
	}
}
func DeactivateProduct(context *gin.Context) {
	fmt.Println("deativating")
	productid := context.Query("id")
	query := "product_id=" + productid

	id := strings.ReplaceAll(productid, "'", "")

	// check if product exist
	productExist, err := FindSingleProduct(id)
	if err != nil {
		response := models.Reply{
			Error:   err.Error(),
			Message: "error finding product",
			Success: false,
		}
		context.JSON(http.StatusBadRequest, response)
		return
	} else if productExist.ProductName == "" {
		response := models.Reply{
			Error:   errors.New("product does not exist").Error(),
			Message: "the product does not exist",
			Success: false,
		}
		context.JSON(http.StatusBadRequest, response)
		return

	} else if productExist.IsDeleted {
		response := models.Reply{
			Error:   errors.New("product deleted").Error(),
			Message: "product is deleted!!Please restore product first",
			Success: false,
		}
		context.JSON(http.StatusBadRequest, response)
		return

	} else if !productExist.IsActive {
		response := models.Reply{
			Data:    productExist,
			Message: "product is not active",
			Success: true,
		}
		context.JSON(http.StatusBadRequest, response)
		return
	} else {
		success, err := DeactivateProductUtil(query)
		if err != nil {
			response := models.Reply{
				Error:   err.Error(),
				Message: "error deactivating  product",
				Success: false,
			}
			context.JSON(http.StatusBadRequest, response)
			return
		} else if !success {
			response := models.Reply{
				Error:   errors.New("could not deactivate product!!try again").Error(),
				Message: "failed in deactivating product",
				Success: false,
			}
			context.JSON(http.StatusBadRequest, response)
			return

		} else {
			response := models.Reply{
				Data:    productExist,
				Message: "succesfully deactivated the product",
				Success: true,
			}
			context.JSON(http.StatusOK, response)
			return

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
		response := models.Reply{
			Error:   err.Error(),
			Message: "error finding product",
			Success: false,
		}
		context.JSON(http.StatusBadRequest, response)
		return
	} else if productExist.ProductName == "" {
		response := models.Reply{
			Data:    Product{},
			Message: "the product does not exist",
			Success: true,
		}
		context.JSON(http.StatusOK, response)
		return
	} else if productExist.IsDeleted {
		response := models.Reply{
			Data:    productExist,
			Message: "product already deleted",
			Success: true,
		}
		context.JSON(http.StatusOK, response)
		return

	} else if productExist.IsActive {
		response := models.Reply{
			Data:    productExist,
			Message: "you cannot delete an active product!! Please deactivate the product first",
			Success: true,
		}
		context.JSON(http.StatusOK, response)
		return
	} else {
		success, err := DeleteProductUtil(query)
		if err != nil {
			response := models.Reply{
				Error:   err.Error(),
				Message: "error deleting  product",
				Success: false,
			}
			context.JSON(http.StatusBadRequest, response)
			return

		} else if !success {
			response := models.Reply{
				Error:   errors.New("could not delete product").Error(),
				Message: "failed in deleting product",
				Success: false,
			}
			context.JSON(http.StatusBadRequest, response)
			return
		} else {
			response := models.Reply{
				Data:    productExist,
				Message: "succesfully deleted the product",
				Success: true,
			}
			context.JSON(http.StatusOK, response)
			return
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
		response := models.Reply{
			Error:   err.Error(),
			Message: "error finding product",
			Success: false,
		}
		context.JSON(http.StatusBadRequest, response)
		return

	} else if productExist.ProductName == "" {
		response := models.Reply{
			Data:    Product{},
			Message: "the product does not exist",
			Success: true,
		}
		context.JSON(http.StatusOK, response)
		return
	} else if !productExist.IsDeleted {
		response := models.Reply{
			Data:    productExist,
			Message: "product is not deleted",
			Success: true,
		}
		context.JSON(http.StatusOK, response)
		return

	} else {
		success, err := RestoreProductUtil(query)
		if err != nil {
			response := models.Reply{
				Error:   err.Error(),
				Message: "error restoring  product",
				Success: false,
			}
			context.JSON(http.StatusBadRequest, response)
			return
		} else if !success {
			response := models.Reply{
				Error:   errors.New("could not restore product!!try again").Error(),
				Message: "failed in restoring product",
				Success: false,
			}
			context.JSON(http.StatusBadRequest, response)
			return
		} else {
			response := models.Reply{
				Data:    productExist,
				Message: "succesfully restoring the product",
				Success: true,
			}
			context.JSON(http.StatusOK, response)
			return

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
		response := models.Reply{
			Error:   err.Error(),
			Message: "error fetching single user products",
			Success: false,
		}
		context.JSON(http.StatusBadRequest, response)
		return

	} else {
		response := models.Reply{
			Data:    products,
			Message: "single user products fetched",
			Success: true,
		}
		context.JSON(http.StatusOK, response)
		return

	}
}

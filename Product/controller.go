package product

import (
	"net/http"
	"time"

	"eleliafrika.com/backend/models"
	"eleliafrika.com/backend/users"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AddProduct(context *gin.Context) {
	var productInput AddProductInput

	if err := context.ShouldBindJSON(&productInput); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	productuuid := uuid.New()
	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02 15:04:05")

	// get current user to add to userid field
	user, err := users.CurrentUser(context)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "could not find user",
		})
	}

	product := models.Product{
		ProductID:          productuuid.String(),
		ProductName:        productInput.ProductName,
		ProductPrice:       productInput.ProductPrice,
		ProductDescription: productInput.ProductDescription,
		UserID:             user.UserID,
		MainImage:          productInput.MainImage,
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
	for _, i := range productInput.ProductImages {

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
	savedProduct, err := product.Save()

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error with save": err.Error(),
			"success":         false,
		})
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"data":    savedProduct,
		"success": true,
		"message": "product has been added succesfully",
	})
}

func GetAllProducts(context *gin.Context) {
	var products []models.Product
	var query string
	var err error
	var data interface{}

	productid := context.Query("id")

	// productid := context.Query("id")
	// productname := context.Query("name")
	// productuser := context.Query("user")
	// productstatus := context.Query("status")
	// producttype := context.Query("type")
	// productbrand := context.Query("brand")
	// productcat := context.Query("cat")
	// productsubcat := context.Query("subcat")
	// productprice := context.Query("price")

	if productid != "" {
		query = "product_id=" + productid

		products, err = Fetchproducts(query)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{
				"error":   err.Error(),
				"success": false,
				"message": "Could not fetch products",
			})
		} else {
			data = gin.H{
				"success":  true,
				"products": products,
			}
		}
	} else {
		query = ""
		products, err = Fetchproducts(query)
		data = gin.H{
			"success":  true,
			"products": products,
		}
	}

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"success": false,
		})
		return
	}

	context.JSON(http.StatusOK, data)
}

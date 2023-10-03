package images

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Getimages(context *gin.Context) {
	productId := context.Param("id")
	images, err := GetSpecificProductImage(productId)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"success": false,
			"message": "Could not fetch products",
		})
	}
	context.JSON(http.StatusOK, gin.H{

		"success": true,
		"message": "Could not fetch products",
		"images":  images,
	})
}

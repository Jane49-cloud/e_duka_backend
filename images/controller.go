package images

import (
	"net/http"

	"eleliafrika.com/backend/models"
	"github.com/gin-gonic/gin"
)

func Getimages(context *gin.Context) {
	productId := context.Param("id")
	images, err := GetSpecificProductImage(productId)
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
		Message: "images fetched succesfully",
		Success: true,
		Data:    images,
	}
	context.JSON(http.StatusOK, response)

}

func UploadOtherImages(imagesString []string, productName string) {

}

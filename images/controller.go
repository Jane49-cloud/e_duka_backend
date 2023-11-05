package images

import (
	"encoding/base64"
	"net/http"
	"strings"

	"eleliafrika.com/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func UploadMainimage(mainImageString string, productName string) (mainimagepath string, err error) {

	imageuuid := uuid.New()

	mainImageFilename := strings.ReplaceAll(productName, " ", "") + imageuuid.String()

	imageBytes, err := base64.StdEncoding.DecodeString(mainImageString)

	if err != nil {
		return "", err
	}

	mainimagepath, err = UploadImageToBucket(productName, "mainimage", imageBytes, mainImageFilename)
	if err != nil {
		return "", err
	}
	return mainimagepath, nil
}

func UploadOtherImages(imagesString []string, productName string) ([]string, error) {
	var imagespath []string

	for _, image := range imagesString {
		imageuuid := uuid.New()

		imagename := strings.ReplaceAll(productName, " ", "") + imageuuid.String()

		imageBytes, err := base64.StdEncoding.DecodeString(image)

		if err != nil {
			return imagespath, err
		}

		imagepath, err := UploadImageToBucket(productName, "other-images", imageBytes, imagename)
		if err != nil {
			return imagespath, err
		}

		imagespath = append(imagespath, imagepath)
	}

	return imagespath, nil
}

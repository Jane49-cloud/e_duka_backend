package images

import (
	"encoding/base64"
	"strings"

	globalutils "eleliafrika.com/backend/global_utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Getimages(context *gin.Context) {
	productId := context.Param("id")
	images, err := GetSpecificProductImage(productId)
	if err != nil {
		globalutils.HandleError("could not the images", err, context)
		return
	}

	globalutils.HandleSuccess("image fetched sucessfully", images, context)

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

package images

import (
	"encoding/base64"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
func UploadMainimage(mainImageString string, productName string) (mainimagepath string, err error) {

	imageuuid := uuid.New()

	mainImageFilename := strings.ReplaceAll(productName, " ", "") + imageuuid.String()

	imageBytes, err := base64.StdEncoding.DecodeString(mainImageString)

	if err != nil {
		return "", err
	}

	mainImagesFolder := "assets/products/"

	if _, err = os.Stat(mainImagesFolder); os.IsNotExist(err) {
		if err = os.Mkdir(mainImagesFolder, 0755); err != nil {
			return "", err
		}
	}
	productFolder := mainImagesFolder + strings.ReplaceAll(productName, " ", "")

	if _, err = os.Stat(productFolder); os.IsNotExist(err) {
		if err = os.Mkdir(productFolder, 0755); err != nil {
			return "", err
		}
	}

	mainImageFolder := productFolder + "/mainimage"

	if _, err = os.Stat(mainImageFolder); os.IsNotExist(err) {
		if err = os.Mkdir(mainImageFolder, 0755); err != nil {
			return "", err
		}
	}

	imagePath := filepath.Join(mainImageFolder, mainImageFilename)

	err = os.WriteFile(imagePath, imageBytes, 0644)
	if err != nil {
		return "", err
	}

	mainimagepath = imagePath

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

		mainImagesFolder := "assets/products/"

		if _, err = os.Stat(mainImagesFolder); os.IsNotExist(err) {
			if err = os.Mkdir(mainImagesFolder, 0755); err != nil {
				return imagespath, err
			}
		}
		productFolder := mainImagesFolder + strings.ReplaceAll(productName, " ", "")

		if _, err = os.Stat(productFolder); os.IsNotExist(err) {
			if err = os.Mkdir(productFolder, 0755); err != nil {
				return imagespath, err
			}
		}

		otherImageFolder := productFolder + "/otherimages"

		if _, err = os.Stat(otherImageFolder); os.IsNotExist(err) {
			if err = os.Mkdir(otherImageFolder, 0755); err != nil {
				return imagespath, err
			}
		}

		imagePath := filepath.Join(otherImageFolder, imagename)

		err = os.WriteFile(imagePath, imageBytes, 0644)
		if err != nil {
			return imagespath, err
		}
		imagespath = append(imagespath, imagePath)
	}

	return imagespath, nil
}

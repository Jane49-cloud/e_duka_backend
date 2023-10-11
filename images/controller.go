package images

import (
	"net/http"
	"os"

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
func UploadMainimage(context *gin.Context, productName string) (mainimagepath string, err error) {
	mainImageFile, err := context.FormFile("mainImage")
	if err != nil {
		return "", err
	}
	mainImageFilename := mainImageFile.Filename
	mainImagesFolder := "assets/mainimages"

	if _, err = os.Stat(mainImagesFolder); os.IsNotExist(err) {
		if err = os.Mkdir(mainImagesFolder, 0755); err != nil {
			return "", err
		}
	}
	mainimagepath = mainImagesFolder + productName + mainImageFilename

	if err = context.SaveUploadedFile(mainImageFile, mainimagepath); err != nil {

		return "", err
	}
	return mainimagepath, nil
}

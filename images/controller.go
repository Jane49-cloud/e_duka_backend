package images

import (
	"bytes"
	"encoding/base64"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
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
func DownloadImageFromBucket(objectKey string) (string, error) {
	awsSecret := os.Getenv("SECRET_KEY")
	awsAccessKey := os.Getenv("ACCESS_KEY")
	token := os.Getenv("TOKEN")

	creds := credentials.NewStaticCredentials(awsAccessKey, awsSecret, token)
	cfg := aws.NewConfig().WithRegion("af-south-1").WithCredentials(creds)

	sess, err := session.NewSession(cfg)
	if err != nil {
		return "", err
	}

	svc := s3.New(sess)
	storageLocation := "e-duka-images"
	input := &s3.GetObjectInput{
		Bucket: aws.String(storageLocation),
		Key:    aws.String(objectKey),
	}

	result, err := svc.GetObject(input)
	if err != nil {
		return "", err
	}
	defer result.Body.Close()

	var imageBuffer bytes.Buffer

	_, err = io.Copy(&imageBuffer, result.Body)

	if err != nil {
		return "", err
	}
	imageString := base64.StdEncoding.EncodeToString(imageBuffer.Bytes())

	return imageString, nil
}

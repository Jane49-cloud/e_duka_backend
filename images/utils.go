package images

import (
	"bytes"
	"fmt"
	"net/http"
	"os"

	"eleliafrika.com/backend/database"
	"eleliafrika.com/backend/models"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func GetProductImages() ([]models.ProductImage, error) {
	var productsImages []models.ProductImage
	err := database.Database.Find(&productsImages).Error
	if err != nil {
		return []models.ProductImage{}, err

	}
	return productsImages, nil
}

func GetSpecificProductImage(productid string) ([]models.ProductImage, error) {
	var productsImages []models.ProductImage
	err := database.Database.Where("product_id=?", productid).Find(&productsImages).Error
	if err != nil {
		return []models.ProductImage{}, err

	}
	return productsImages, nil
}

func UploadImageToBucket(productName string, imagefolder string, imageBytes []byte, imagename string) (string, error) {

	awsSecret := os.Getenv("SECRET_KEY")
	awsAccessKey := os.Getenv("ACCESS_KEY")
	token := os.Getenv("TOKEN")

	creds := credentials.NewStaticCredentials(awsAccessKey, awsSecret, token)
	cfg := aws.NewConfig().WithRegion("af-south-1").WithCredentials(creds)

	sess, err := session.NewSession(cfg)
	if err != nil {
		return "", err
	}

	reader := bytes.NewReader(imageBytes)

	fileSize := int64(len(imageBytes))
	fileType := http.DetectContentType(imageBytes)

	svc := s3.New(sess)
	objectKey := "assets/productImages/" + productName + "/" + imagefolder + "/" + imagename
	storageLocation := "e-duka-images"

	input := &s3.PutObjectInput{
		Body:          reader,
		Bucket:        aws.String(storageLocation),
		Key:           aws.String(objectKey),
		ContentType:   aws.String(fileType),
		ContentLength: aws.Int64(fileSize),
	}

	_, err = svc.PutObject(input)
	if err != nil {
		return "", err
	}
	fmt.Printf("main image string \n%s\n", objectKey)
	return objectKey, nil

}

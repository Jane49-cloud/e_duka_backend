package images

import (
	"bytes"
	"encoding/base64"
	"io"
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

func GetSpecificProductImage(productid string) ([]string, error) {
	var productsImages []models.ProductImage
	var images []string
	err := database.Database.Where("product_id=?", productid).Find(&productsImages).Error
	if err != nil {
		return []string{}, err
	}
	for image := range productsImages {
		images = append(images, productsImages[image].ImageUrl)
	}
	return images, nil
}

func UploadImageToBucket(name string, imagefolder string, imageBytes []byte, imagename string) (string, error) {

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
	objectKey := "assets/productImages/" + name + "/" + imagefolder + "/" + imagename
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

	return objectKey, nil

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

func DeleteImageFromBucket(bucketName string, objectKey string) (bool, error) {
	awsSecret := os.Getenv("SECRET_KEY")
	awsAccessKey := os.Getenv("ACCESS_KEY")
	token := os.Getenv("TOKEN")

	creds := credentials.NewStaticCredentials(awsAccessKey, awsSecret, token)
	cfg := aws.NewConfig().WithRegion("af-south-1").WithCredentials(creds)

	sess, err := session.NewSession(cfg)
	if err != nil {
		return false, err
	}

	svc := s3.New(sess)

	input := &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}

	_, err = svc.DeleteObject(input)
	if err != nil {
		return false, err
	}

	return true, nil
}

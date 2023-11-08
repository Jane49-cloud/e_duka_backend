package images

import (
	"context"
	"encoding/base64"
	"log"

	"eleliafrika.com/backend/database"
	"eleliafrika.com/backend/models"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
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
func UploadImageToFireBase(imageString string) (string, error) {
	config := &firebase.Config{
		StorageBucket: "eduka-f19e5.appspot.com",
	}

	opt := option.WithCredentialsFile("./key.json")
	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Storage(context.Background())
	if err != nil {
		log.Fatalln(err)
	}

	bucket, err := client.DefaultBucket()
	if err != nil {
		log.Fatalln(err)
	}

	decodedImage, err := base64.StdEncoding.DecodeString(imageString)

	if err != nil {
		return "", err
	}

	// Generate a unique image name or use an existing one
	imageName := "newimage.jpg"

	object := bucket.Object("images/" + imageName)
	wc := object.NewWriter(context.Background())
	wc.ContentType = "image/jpeg"

	_, err = wc.Write(decodedImage)
	if err != nil {
		return "", err
	}

	if err := wc.Close(); err != nil {
		return "", err
	}

	// Get the download URL
	downloadURL, err := object.Attrs(context.Background())
	if err != nil {
		return "", err
	}

	return downloadURL.MediaLink, nil

}

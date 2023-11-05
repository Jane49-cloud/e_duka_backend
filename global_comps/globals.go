package globalcomps

import (
	"fmt"
	"log"

	"eleliafrika.com/backend/admin"
	"eleliafrika.com/backend/brands"
	"eleliafrika.com/backend/category"
	"eleliafrika.com/backend/comments"
	"eleliafrika.com/backend/database"
	"eleliafrika.com/backend/images"
	"eleliafrika.com/backend/mainad"
	"eleliafrika.com/backend/product"
	subcategory "eleliafrika.com/backend/subcategories"
	"eleliafrika.com/backend/users"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func LoadDatabase() {
	database.Connect()
	// database.Database.AutoMigrate(&admin.SystemAdmin{}, &users.User{}, &models.Brand{}, &models.Category{}, &models.SubCategory{}, &models.Comment{}, &product.Product{})
	// database.Database.AutoMigrate(&models.ProductImage{}, &product.Product{})
}

func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
		log.Fatal(err)
	}
}

func ServeApplication() {
	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Header("Referrer-Policy", "no-referrer")
		c.Next()
	})
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"POST", "GET", "PUT"}
	config.AllowHeaders = []string{"Content-Type", "x-access-token"}
	router.Use(cors.New(config))

	users.UserRoutes(router)
	product.ProductRoutes(router)
	images.Imagesroutes(router)
	comments.Commentroutes(router)
	category.CategoryRoutes(router)
	subcategory.SubCategoryRoutes(router)
	brands.BrandRoutes(router)
	mainad.Mainadsroutes(router)
	admin.AdminRoutes(router)

	// Load the SSL certificate and key
	certFile := "./fullchain.pem" // Update this with the path to your certificate file
	keyFile := "./privkey.pem"    // Update this with the path to your private key file

	// Run the server with TLS/HTTPS
	//	if err := router.Run(":8000");
	if err := router.RunTLS("192.0.0.105:8000", certFile, keyFile); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	fmt.Println("Server running on port 8000 (HTTPS)")
}

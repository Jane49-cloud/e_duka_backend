package globalcomps

import (
	"fmt"
	"log"

	product "eleliafrika.com/backend/Product"
	"eleliafrika.com/backend/brands"
	"eleliafrika.com/backend/category"
	"eleliafrika.com/backend/comments"
	"eleliafrika.com/backend/database"
	"eleliafrika.com/backend/images"
	"eleliafrika.com/backend/mainad"
	"eleliafrika.com/backend/models"
	subcategory "eleliafrika.com/backend/subcategories"
	"eleliafrika.com/backend/users"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func LoadDatabase() {
	database.Connect()
	// database.Database.AutoMigrate(&models.User{}, &models.Brand{}, &models.Category{}, &models.SubCategory{}, &models.Comment{}, &models.Product{})
	database.Database.AutoMigrate(&models.User{})
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

	users.AuthRoutes(router)
	product.ProductRoutes(router)
	images.Imagesroutes(router)
	comments.Commentroutes(router)
	category.CategoryRoutes(router)
	subcategory.SubCategoryRoutes(router)
	brands.BrandRoutes(router)
	mainad.Mainadsroutes(router)

	router.Run(":8000")
	fmt.Println("Server running on port 8000")
}

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
	subcategory "eleliafrika.com/backend/subcategories"
	"eleliafrika.com/backend/users"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func LoadDatabase() {
	database.Connect()
	// database.Database.AutoMigrate(&models.User{}, &models.Brand{}, &models.Category{}, &models.SubCategory{}, &models.Comment{}, &models.Product{})
	// database.Database.AutoMigrate(&models.Comment{})
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

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true

	router.Use(cors.New(config))

	users.AuthRoutes(router)
	product.PostRoutes(router)
	images.Imagesroutes(router)
	comments.Commentroutes(router)
	category.CategoryRoutes(router)
	subcategory.SubCategoryRoutes(router)
	brands.BrandRoutes(router)

	router.Run(":8000")
	fmt.Println("Server running on port 8000")
}

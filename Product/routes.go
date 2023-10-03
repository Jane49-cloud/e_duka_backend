package product

import (
	"eleliafrika.com/backend/images"
	"eleliafrika.com/backend/users"
	"github.com/gin-gonic/gin"
)

func PostRoutes(router *gin.Engine) {
	authRoutes := router.Group("/products")

	{
		authRoutes.POST("/addproduct", users.JWTAuthMiddleWare(), AddProduct)
		authRoutes.GET("/getproducts", GetAllProducts)
		authRoutes.GET("/getproducts/single/", GetAllProducts)
		authRoutes.GET("/getimages", images.Getimages)
	}
}

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
		authRoutes.GET("/getproducts/single/", GetSingleProduct)
		authRoutes.POST("/updateproduct/", UpdateProduct)
		authRoutes.GET("/getimages", images.Getimages)
		authRoutes.POST("/deleteproduct", DeleteProduct)
		authRoutes.POST("/restore", RestoreProduct)
		authRoutes.POST("/activate", ActivateProduct)
		authRoutes.POST("/deactivate", DeactivateProduct)
	}
}

package product

import (
	"eleliafrika.com/backend/users"
	"github.com/gin-gonic/gin"
)

func ProductRoutes(router *gin.Engine) {
	authRoutes := router.Group("/products")

	{
		authRoutes.POST("/addproduct", users.JWTAuthMiddleWare(), AddProduct)
		authRoutes.GET("/getproducts", GetAllProducts)
		authRoutes.GET("/getads", GetAllAds)
		authRoutes.GET("/getproducts/single/:id", GetSingleProduct)
		authRoutes.GET("/getads/single/:id", GetSingleAd)
		authRoutes.POST("/updateproduct/", UpdateProduct)
		authRoutes.POST("/deleteproduct", users.JWTAuthMiddleWare(), DeleteProduct)
		authRoutes.POST("/restore", users.JWTAuthMiddleWare(), RestoreProduct)
		authRoutes.POST("/activate", users.JWTAuthMiddleWare(), ActivateProduct)
		authRoutes.POST("/approve", users.JWTAuthMiddleWare(), ApproveProduct)
		authRoutes.POST("/deactivate", users.JWTAuthMiddleWare(), DeactivateProduct)
	}
}

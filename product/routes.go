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
		authRoutes.GET("/getproducts/singleuserproduct", users.JWTAuthMiddleWare(), FetchSingleUserProducts)
		authRoutes.GET("/getads/singleuserads", FetchSingleUserAds)
		authRoutes.POST("/deleteproduct", users.JWTAuthMiddleWare(), DeleteProduct)
		authRoutes.POST("/restore", users.JWTAuthMiddleWare(), RestoreProduct)
		authRoutes.POST("/activate", users.JWTAuthMiddleWare(), ActivateProduct)
		authRoutes.POST("/deactivate", users.JWTAuthMiddleWare(), DeactivateProduct)
	}
}

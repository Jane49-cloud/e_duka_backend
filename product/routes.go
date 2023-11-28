package product

import (
	"eleliafrika.com/backend/users"
	"github.com/gin-gonic/gin"
)

func ProductRoutes(router *gin.Engine) {
	productRoutes := router.Group("/products")

	{
		productRoutes.POST("/addproduct", users.JWTAuthMiddleWare(), AddProduct)
		productRoutes.GET("/getproducts", GetAllProducts)
		productRoutes.GET("/getproductsdata", GetAllAds)
		productRoutes.GET("/getproducts/single/:id", GetSingleProduct)
		productRoutes.GET("/getads/single/:id", GetSingleAd)
		productRoutes.POST("/updateproduct", UpdateProduct)
		productRoutes.GET("/getproducts/singleuserproduct", users.JWTAuthMiddleWare(), FetchSingleUserProducts)
		productRoutes.GET("/getads/singleuserads", FetchSingleUserAds)
		productRoutes.POST("/deleteproduct", users.JWTAuthMiddleWare(), DeleteProduct)
		productRoutes.POST("/restore", users.JWTAuthMiddleWare(), RestoreProduct)
		productRoutes.POST("/activate", users.JWTAuthMiddleWare(), ActivateProduct)
		productRoutes.POST("/deactivate", users.JWTAuthMiddleWare(), DeactivateProduct)
	}
}

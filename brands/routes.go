package brands

import (
	"eleliafrika.com/backend/users"
	"github.com/gin-gonic/gin"
)

func BrandRoutes(router *gin.Engine) {
	brandroutes := router.Group("/brands")
	{
		brandroutes.POST("/addbrand", users.JWTAuthMiddleWare(), AddBrand)
		brandroutes.GET("/getbrands", GetAllBrands)
		brandroutes.POST("/delete/:name", DeleteBrand)
	}
}

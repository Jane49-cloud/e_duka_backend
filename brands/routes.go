package brands

import (
	"eleliafrika.com/backend/users"
	"github.com/gin-gonic/gin"
)

func BranRoutes(router *gin.Engine) {
	brandroutes := router.Group("/brands")
	{
		brandroutes.POST("/addbrand", users.JWTAuthMiddleWare(), AddBrand)
		brandroutes.GET("/getbrands", GetAllBrands)
	}
}

package category

import (
	"eleliafrika.com/backend/users"
	"github.com/gin-gonic/gin"
)

func CategoryRoutes(router *gin.Engine) {
	categoryRoutes := router.Group("/category")
	{
		categoryRoutes.POST("/addcategory", users.JWTAuthMiddleWare(), CreateCategory)
		categoryRoutes.GET("/getcategories", GetCategories)
	}
}

package category

import (
	"eleliafrika.com/backend/users"
	"github.com/gin-gonic/gin"
)

func CategoryRoutes(router *gin.Engine) {
	categoryRoutes := router.Group("/categories")
	{
		categoryRoutes.POST("/addcategory", users.JWTAuthMiddleWare(), CreateCategory)
		categoryRoutes.GET("/getcategories", GetCategories)
		categoryRoutes.POST("/delete/:name", users.JWTAuthMiddleWare(), DeleteCategory)
		// categoryRoutes.POST("/delete/:name", users.JWTAuthMiddleWare(), DeleteCategory)
	}
}

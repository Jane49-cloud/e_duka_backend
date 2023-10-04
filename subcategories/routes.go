package subcategory

import (
	"eleliafrika.com/backend/users"
	"github.com/gin-gonic/gin"
)

func SubCategoryRoutes(router *gin.Engine) {
	categoryRoutes := router.Group("/subcategories")
	{
		categoryRoutes.POST("/addsubcategory", users.JWTAuthMiddleWare(), CreateSubCategory)
		categoryRoutes.GET("/getsubcategories", GetCategories)
		categoryRoutes.POST("/delete/:name", users.JWTAuthMiddleWare(), DeleteCategory)
	}
}

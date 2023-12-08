package packages

import (
	"eleliafrika.com/backend/users"
	"github.com/gin-gonic/gin"
)

func PackagesRoutes(router *gin.Engine) {
	packagesRoutes := router.Group("/packages", users.JWTAuthMiddleWare())
	{
		packagesRoutes.POST("/create", CreatePackage)
		packagesRoutes.POST("/update", UpdatePackage)
		packagesRoutes.GET("/getsinglepackage", FetchSinglePackage)
		packagesRoutes.GET("/getallpackages", FetchAllPackages)

	}
}

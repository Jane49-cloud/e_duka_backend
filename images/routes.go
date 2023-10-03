package images

import (
	"github.com/gin-gonic/gin"
)

func Imagesroutes(router *gin.Engine) {
	imagesRoutes := router.Group("/images")

	{
		imagesRoutes.GET("/getimages/:id", Getimages)
	}
}

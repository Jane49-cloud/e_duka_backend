package comments

import (
	"eleliafrika.com/backend/users"
	"github.com/gin-gonic/gin"
)

func Commentroutes(router *gin.Engine) {
	commentsRoutes := router.Group("/comments")
	commentsRoutes.Use(users.JWTAuthMiddleWare())
	{
		commentsRoutes.POST("/create", MakeComment)
		commentsRoutes.GET("/get/:id", GetComments)
		commentsRoutes.POST("/delete/:id", DeleteComment)
	}
}

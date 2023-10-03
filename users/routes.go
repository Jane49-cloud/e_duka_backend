package users

import "github.com/gin-gonic/gin"

func AuthRoutes(router *gin.Engine) {
	authRoutes := router.Group("/user/auth")
	{
		authRoutes.POST("/signup", Register)
		authRoutes.POST("/signin", Login)
	}
}

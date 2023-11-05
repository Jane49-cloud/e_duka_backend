package admin

import (
	"eleliafrika.com/backend/users"
	"github.com/gin-gonic/gin"
)

func AdminRoutes(router *gin.Engine) {
	authRoutes := router.Group("/admin")
	{

		authRoutes.POST("/register", Register)
		authRoutes.POST("/login", Login)
		authRoutes.POST("/approveuser", users.JWTAuthMiddleWare(), ApproveUser)
		authRoutes.POST("/revokeuser", users.JWTAuthMiddleWare(), RevokeUser)
		authRoutes.GET("/fetchusers", users.JWTAuthMiddleWare(), FetchSellers)
		authRoutes.POST("/approveproduct", users.JWTAuthMiddleWare(), ApproveProduct)
	}
}
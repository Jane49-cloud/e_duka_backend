package admin

import (
	"eleliafrika.com/backend/users"
	"github.com/gin-gonic/gin"
)

func AdminRoutes(router *gin.Engine) {
	authRoutes := router.Group("/admin")
	{

		authRoutes.POST("/approveuser", users.JWTAuthMiddleWare(), ApproveUser)
		authRoutes.POST("/revokeuser", users.JWTAuthMiddleWare(), RevokeUser)
		authRoutes.GET("/fetchusers", users.JWTAuthMiddleWare(), FetchSellers)
		authRoutes.POST("/approve", users.JWTAuthMiddleWare(), ApproveProduct)
	}
}

package chat

import (
	"eleliafrika.com/backend/users"
	"github.com/gin-gonic/gin"
)

func ChatRoutes(router *gin.Engine) {
	requestRoutes := router.Group("/chats", users.JWTAuthMiddleWare())
	{
		requestRoutes.POST("/sendmessage", SendMessage)
		requestRoutes.GET("/getconversations", GetMessages)
		requestRoutes.GET("/getsinglechat", GetMessages)
		requestRoutes.POST("/deletechat", DeleteRequest)
	}
}

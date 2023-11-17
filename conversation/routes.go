package conversation

import (
	"eleliafrika.com/backend/users"
	"github.com/gin-gonic/gin"
)

func ConversationRoutes(router *gin.Engine) {
	conversationRoutes := router.Group("/conversation", users.JWTAuthMiddleWare())
	{
		conversationRoutes.POST("/startconversation", StartNewConversation)
		conversationRoutes.POST("/getsellerconversations", GetSellerConversations)
		conversationRoutes.POST("/getclientconversations", GetClientConversations)
	}
}

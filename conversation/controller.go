package conversation

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"eleliafrika.com/backend/chat"
	"eleliafrika.com/backend/models"
	"eleliafrika.com/backend/users"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func StartNewConversation(context *gin.Context) {

	var conversationInput ConversationInput

	if err := context.ShouldBindJSON(&conversationInput); err != nil {
		response := models.Reply{
			Error:   err.Error(),
			Message: "wrong data format from the user",
			Success: false,
		}
		context.JSON(http.StatusBadRequest, response)
		return
	}
	conversationuuid := uuid.New()
	chatuuid := uuid.New()
	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02 15:04:05")

	user, err := users.CurrentUser(context)
	if err != nil {
		response := models.Reply{
			Message: "error fetching user",
			Success: false,
			Error:   err.Error(),
		}
		context.JSON(http.StatusUnauthorized, response)
		return
	} else if user.Firstname == "" {
		response := models.Reply{
			Message: "user not found",
			Success: false,
			Error:   errors.New("user not found").Error(),
		}
		context.JSON(http.StatusUnauthorized, response)
		return
	} else {
		conversation := Conversation{
			ConversationId: conversationuuid.String(),
			CustomerId:     user.UserID,
			SellerId:       conversationInput.SellerId,
			LastMessage:    conversationInput.Message,
			MessagesCount:  1,
			IsViewed:       false,
			DateAdded:      formattedTime,
		}

		_, err = conversation.Save()
		if err != nil {
			response := models.Reply{
				Message: "could not create conversation",
				Error:   err.Error(),
				Success: false,
			}
			context.JSON(http.StatusBadRequest, response)
			return
		} else {
			chat := chat.Chat{
				ChatID:         chatuuid.String(),
				ConversationId: conversationuuid.String(),
				SenderID:       user.UserID,
				ReceiverId:     conversation.SellerId,
				Message:        conversationInput.Message,
				TimeSent:       formattedTime,
				IsViewed:       false,
			}

			_, err := chat.Save()
			if err != nil {
				respose := models.Reply{
					Message: "could not start conversation",
					Error:   err.Error(),
					Success: false,
				}
				context.JSON(http.StatusBadRequest, respose)
				return
			} else {

				response := models.Reply{
					Message: "conversation created",
					Success: true,
				}
				context.JSON(http.StatusCreated, response)
				return
			}
		}
	}
}
func GetClientConversations(context *gin.Context) {
	user, err := users.CurrentUser(context)

	if err != nil {
		response := models.Reply{
			Error:   err.Error(),
			Message: "error authorizing user",
			Success: false,
		}
		context.JSON(http.StatusUnauthorized, response)
		return
	} else if user.Firstname == "" {
		response := models.Reply{
			Error:   errors.New("user does not exist").Error(),
			Message: "user does not exist",
			Success: false,
		}
		context.JSON(http.StatusUnauthorized, response)
		return
	}

	conversations, err := FindClientConversation(strings.ReplaceAll(user.UserID, " ", ""))
	if err != nil {
		response := models.Reply{
			Error:   err.Error(),
			Message: "could not fetch user conversations",
			Success: false,
		}
		context.JSON(http.StatusBadRequest, response)
		return
	} else {
		response := models.Reply{
			Data:    conversations,
			Message: "conversations fetched",
			Success: true,
		}
		context.JSON(http.StatusOK, response)
		return
	}
}
func GetSellerConversations(context *gin.Context) {
	user, err := users.CurrentUser(context)

	if err != nil {
		response := models.Reply{
			Error:   err.Error(),
			Message: "error authorizing user",
			Success: false,
		}
		context.JSON(http.StatusUnauthorized, response)
		return
	} else if user.Firstname == "" {
		response := models.Reply{
			Error:   errors.New("user does not exist").Error(),
			Message: "user does not exist",
			Success: false,
		}
		context.JSON(http.StatusUnauthorized, response)
		return
	}

	conversations, err := FindSellerConversation(strings.ReplaceAll(user.UserID, " ", ""))
	if err != nil {
		response := models.Reply{
			Error:   err.Error(),
			Message: "could not fetch user conversations",
			Success: false,
		}
		context.JSON(http.StatusBadRequest, response)
		return
	} else {
		response := models.Reply{
			Data:    conversations,
			Message: "conversations fetched",
			Success: true,
		}
		context.JSON(http.StatusOK, response)
		return
	}
}

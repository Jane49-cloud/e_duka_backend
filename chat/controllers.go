package chat

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"eleliafrika.com/backend/models"
	"eleliafrika.com/backend/users"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func SendMessage(context *gin.Context) {
	var chatInput ChatInput

	conversationid := context.Query("id")

	if err := context.ShouldBindJSON(&chatInput); err != nil {
		response := models.Reply{
			Message: "could not bind data",
			Error:   errors.New("data from user not in correct format").Error(),
			Success: false,
		}
		context.JSON(http.StatusBadRequest, response)
		return
	}

	chatuuid := uuid.New()
	currentTime := time.Now()
	formattedTime := currentTime.Format("2000-01-02 15:04:00")

	// get the current user
	user, err := users.CurrentUser(context)

	if err != nil {
		response := models.Reply{
			Message: "error authoriizing user",
			Success: false,
			Error:   err.Error(),
		}
		context.JSON(http.StatusUnauthorized, response)
		return
	}

	// check is user exists
	if user.Firstname == "" {
		response := models.Reply{
			Message: "user does not exist",
			Success: false,
			Error:   errors.New("could not find the user").Error(),
		}
		context.JSON(http.StatusUnauthorized, response)
		return
	}

	chat := Chat{
		ChatID:         chatuuid.String(),
		ConversationId: strings.ReplaceAll(conversationid, "'", ""),
		SenderID:       user.UserID,
		ReceiverId:     chatInput.ReceiverId,
		Message:        chatInput.Message,
		TimeSent:       formattedTime,
		IsViewed:       false,
	}

	_, err = chat.Save()

	if err != nil {
		response := models.Reply{
			Message: "could not make request",
			Success: false,
			Error:   err.Error(),
		}
		context.JSON(http.StatusBadRequest, response)
		return
	}

	response := models.Reply{
		Message: "text sent successfully",
		Success: true,
	}
	context.JSON(http.StatusOK, response)
}

func GetMessages(context *gin.Context) {

	conversationid := context.Query("id")

	user, err := users.CurrentUser(context)

	if err != nil {
		response := models.Reply{
			Error:   err.Error(),
			Message: "error authorizing user",
			Success: false,
		}
		context.JSON(http.StatusUnauthorized, response)
		return
	}
	if user.Firstname == "" {
		response := models.Reply{
			Error:   errors.New("user does not exist").Error(),
			Message: "user does not exist",
			Success: false,
		}
		context.JSON(http.StatusUnauthorized, response)
		return
	}

	chats, err := GetChats(strings.ReplaceAll(conversationid, "'", ""))
	if err != nil {
		response := models.Reply{
			Message: "error fetching chats",
			Error:   err.Error(),
			Success: false,
		}
		context.JSON(http.StatusBadRequest, response)
		return
	}

	response := models.Reply{
		Message: "chats fetched for " + user.Firstname,
		Data:    chats,
		Success: true,
	}

	context.JSON(http.StatusOK, response)
}

func DeleteRequest(context *gin.Context) {
	response := models.Reply{
		Message: "deleting single request",
		Data:    "deleting the single request",
		Success: true,
	}
	context.JSON(http.StatusOK, response)
}

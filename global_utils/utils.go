package globalutils

import (
	"net/http"

	"eleliafrika.com/backend/models"
	"github.com/gin-gonic/gin"
)

func HandleError(message string, err error, context *gin.Context) {
	response := models.Reply{
		Message: message,
		Error:   err.Error(),
		Success: false,
	}
	context.JSON(http.StatusBadRequest, response)
	return
}

func HandleSuccess(message string, data interface{}, context *gin.Context) {
	response := models.Reply{
		Message: message,
		Data:    data,
		Success: true,
	}
	context.JSON(http.StatusOK, response)
	return
}

package globalutils

import (
	"net/http"

	"eleliafrika.com/backend/models"
	"github.com/gin-gonic/gin"
)

func HandleError(message string, err error, context *gin.Context) {
	response := gin.H{
		"Message": message,
		"Error":   err.Error(),
		"Success": false,
	}
	context.JSON(http.StatusBadRequest, response)

}

func HandleSuccess(message string, data interface{}, context *gin.Context) {
	response := models.Reply{
		Message: message,
		Data:    data,
		Success: true,
	}
	context.JSON(http.StatusOK, response)

}
func UnAuthenticated(context *gin.Context) {
	response := gin.H{
		"Message": "not authenticated",
		"Success": false,
	}
	context.JSON(http.StatusNetworkAuthenticationRequired, response)

}

func UnAuthorized(context *gin.Context) {
	response := gin.H{
		"Message": "not authorized for the action",
		"Success": false,
	}
	context.JSON(http.StatusUnauthorized, response)

}

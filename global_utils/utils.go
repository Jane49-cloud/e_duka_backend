package globalutils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

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

package globalutils

import (
	"context"
	"log"
	"net/http"

	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
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

func InitFirebaseApp() (*firebase.App, error) {
	ctx := context.Background()

	// Set up Firebase Admin SDK credentials.
	opt := option.WithCredentialsFile("./key.json")

	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalf("Error initializing Firebase app: %v\n", err)
		return nil, err
	}

	return app, nil
}

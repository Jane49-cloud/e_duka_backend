package globalutils

import (
	"context"
	"log"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

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

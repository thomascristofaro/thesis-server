package utility

import (
	"context"
	"encoding/base64"
	"errors"
	"os"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
)

func getDecodedFireBaseKey() ([]byte, error) {
	FirebaseKey := os.Getenv("FIREBASE_KEY")
	if FirebaseKey == "" {
		return nil, errors.New("FirebaseKey not found")
	}
	return base64.StdEncoding.DecodeString(FirebaseKey)
}

func SendFirebaseNotification(context context.Context, deviceToken string,
	title string, body string, data map[string]string) error {

	decodedKey, err := getDecodedFireBaseKey()
	if err != nil {
		return err
	}

	opt := option.WithCredentialsJSON(decodedKey)
	app, err := firebase.NewApp(context, nil, opt)
	if err != nil {
		return err
	}

	fcmClient, err := app.Messaging(context)
	if err != nil {
		return err
	}

	_, err = fcmClient.Send(context, &messaging.Message{
		Data: data,
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Token: deviceToken,
	})
	if err != nil {
		return err
	}

	return nil
}

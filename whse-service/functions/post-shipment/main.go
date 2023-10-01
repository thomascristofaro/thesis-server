package main

import (
	"context"
	"encoding/json"
	"fmt"
	"thesis/lib/utility"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(ctx context.Context, snsEvent events.SNSEvent) error {
	message, err := utility.ConvertSNSEventToMessage(snsEvent)
	if err != nil {
		return err
	}
	fmt.Printf("Message: %s", message.Body)

	// registra spedizione
	var v map[string]interface{}
	json.Unmarshal(message.Body, &v)
	v["event"] = "ChangeOrderStatusFromShpt"
	message.Body, err = json.Marshal(v)
	if err != nil {
		return err
	}

	// cambia stato ordine
	err = utility.SendSQSMessage(ctx,
		"ChangeOrderStatusQueueUrl",
		message, true)

	if err != nil {
		return err
	}

	v["event"] = "OnAfterPostingShpt"
	message.Body, err = json.Marshal(v)
	if err != nil {
		return err
	}

	// chiama evento post spedizione
	err = utility.SendSNSMessage(ctx,
		"OnAfterPostingShptTopicArn",
		message)

	if err != nil {
		return err
	}

	return nil
}

func main() {
	lambda.Start(Handler)
}

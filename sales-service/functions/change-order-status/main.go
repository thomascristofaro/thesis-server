package main

import (
	"context"
	"fmt"
	"thesis/lib/utility"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	message, err := utility.ConvertSQSEventToMessage(sqsEvent)
	if err != nil {
		return err
	}
	fmt.Printf("Message: %s", message.Body)

	return nil
}

func main() {
	lambda.Start(Handler)
}

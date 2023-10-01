package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	for _, record := range sqsEvent.Records {
		fmt.Printf("[%s] Message = %s \n", record.EventSource, record.Body)
		b, err := json.Marshal(record)
		if err != nil {
			return err
		}
		fmt.Printf("[%s] JSON MSG = %s \n", record.EventSource, string(b))
	}
	return nil
}

func main() {
	lambda.Start(Handler)
}

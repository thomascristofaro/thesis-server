package main

import (
	"context"
	"fmt"
	"thesis/lib/utility"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(ctx context.Context, snsEvent events.SNSEvent) error {
	for _, record := range snsEvent.Records {
		snsRecord := record.SNS
		fmt.Printf("[%s %s] Message = %s \n", record.EventSource, snsRecord.Timestamp, snsRecord.Message)
	}

	// registra spedizione

	// cambia stato ordine
	err := utility.SendSQSMessage(ctx,
		"ChangeOrderStatusQueueUrl",
		[]byte("ChangeOrderStatus Message!"),
		map[string]string{
			"language":   "en",
			"importance": "high",
		})

	if err != nil {
		return err
	}

	// chiama evento post spedizione
	err = utility.SendSNSMessage(ctx,
		"OnAfterPostingShptTopicArn",
		[]byte("OnAfterPostingShpt Message!"),
		map[string]string{
			"language":   "en",
			"importance": "high",
		})

	if err != nil {
		return err
	}

	return nil
}

func main() {
	lambda.Start(Handler)
}

package main

import (
	"context"
	"encoding/json"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"gocloud.dev/pubsub"
	_ "gocloud.dev/pubsub/awssnssqs"
)

type Response events.APIGatewayProxyResponse

func Handler(ctx context.Context) (Response, error) {
	// Registro l'ordine

	// lancio l'evento di registrazione ordine
	topicARN := os.Getenv("OnAfterPostingOrderTopicArn")
	topic, err := pubsub.OpenTopic(ctx, "awssns:///"+topicARN+"?region=us-east-1")
	if err != nil {
		return Response{}, err
	}
	defer topic.Shutdown(ctx)

	err = topic.Send(ctx, &pubsub.Message{
		Body: []byte("Hello, World!\n"),
		// Metadata is optional and can be nil.
		Metadata: map[string]string{
			"language":   "en",
			"importance": "high",
		},
	})
	if err != nil {
		return Response{}, err
	}

	// risposta HTTP
	body, err := json.Marshal(map[string]interface{}{
		"message": "Go Serverless v1.0! Your function executed successfully!",
	})
	if err != nil {
		return Response{StatusCode: 404}, err
	}

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            string(body),
		Headers: map[string]string{
			"Content-Type":           "application/json",
			"X-MyCompany-Func-Reply": "hello-handler",
		},
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}

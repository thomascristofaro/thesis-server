package main

import (
	"context"
	"encoding/json"
	"thesis/lib/utility"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Response events.APIGatewayProxyResponse

func Handler(ctx context.Context) (Response, error) {
	// Registro l'ordine

	err := utility.SendSNSMessage(ctx,
		"OnAfterPostingOrderTopicArn",
		[]byte("OnAfterPostingOrder Message!"),
		map[string]string{
			"language":   "en",
			"importance": "high",
		})

	if err != nil {
		return Response{StatusCode: 500}, err
	}

	// risposta HTTP
	body, err := json.Marshal(map[string]interface{}{
		"message": "Go Serverless v1.0! Your function executed successfully!",
	})
	if err != nil {
		return Response{StatusCode: 500}, err
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

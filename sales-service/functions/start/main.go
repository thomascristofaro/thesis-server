package main

import (
	"context"
	"encoding/json"
	"strconv"
	"thesis/lib/utility"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Response events.APIGatewayProxyResponse

func Handler(ctx context.Context) (Response, error) {
	// Registro l'ordine

	for i := 0; i < 10; i++ {
		body, err := json.Marshal(map[string]interface{}{
			"order_id": "OR" + strconv.Itoa(i),
			"event":    "OnAfterPostingOrder",
		})
		if err != nil {
			return Response{StatusCode: 500}, err
		}
		message := utility.Message{
			Body: body,
			Metadata: map[string]string{
				"device_id": "XXXXXXXX",
			},
		}

		err = utility.SendSNSMessage(ctx,
			"OnAfterPostingOrderTopicArn",
			message)

		if err != nil {
			return Response{StatusCode: 500}, err
		}
	}

	// risposta HTTP
	body, err := json.Marshal(map[string]interface{}{
		"message": "SNS Message Sent!",
	})
	if err != nil {
		return Response{StatusCode: 500}, err
	}

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            string(body),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}

package main

import (
	"context"
	"encoding/json"
	"thesis/lib/database"
	"thesis/lib/utility"
	"thesis/whse-service/models"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// TODO Log error
func PostShipment(ctx context.Context, message utility.Message) error {
	device_id := message.Metadata["device_id"]

	// registra spedizione
	body := map[string]interface{}{}
	err := json.Unmarshal(message.Body, &body)
	if err != nil {
		return err
	}

	m := database.NewModel(models.NewShipmentHeader())
	m.Open()
	defer m.Close()
	model := m.Model.(*models.ShipmentHeader)
	model.No = body["No"].(string)
	if !m.Create() {
		return m.GetLastError()
	}

	// cambia stato ordine
	var bodyBuff []byte
	bodyBuff, err = json.Marshal(map[string]string{
		"response": "OK",
	})
	if err != nil {
		return err
	}
	metadata := map[string]string{
		"device_id": device_id,
		"function":  "PostShipment",
		"service":   "SQS",
		"event":     "OnChangeOrderStatus",
	}
	err = utility.SendSQSMessage(ctx,
		"ChangeOrderStatusQueueUrl",
		utility.Message{
			Body:     bodyBuff,
			Metadata: metadata,
		}, true)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	lambda.Start(func(ctx context.Context, snsEvent events.SNSEvent) error {
		message, err := utility.ConvertSNSEventToMessage(snsEvent)
		if err != nil {
			return err
		}
		return PostShipment(ctx, message)
	})

	// PostShipment(context.Background(), utility.Message{
	// 	Body: []byte(`{"No":"OR0001"}`),
	// 	Metadata: map[string]string{
	// 		"device_id": "eKtDSgRajJNmpcgKSQDTod:APA91bEE-t0pxZpNS5uG4jQEkV0I0P58fBBavNf9MldWxVp8xQJunv5UR7tReQ_seAWK1IxsdrinYANyqies47tmKpNStemyFTccwZiJJ8itsACJBQYlMVTw_rfG5ZO-_iIdgX5Y_QEZ",
	// 	},
	// })
}

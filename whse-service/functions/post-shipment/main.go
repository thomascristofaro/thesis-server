package main

import (
	"context"
	"encoding/json"
	"errors"
	"thesis/lib/database"
	"thesis/lib/utility"
	"thesis/whse-service/models"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func PostShipment(ctx context.Context, message *utility.Message) error {
	utility.BuildLogMetadata("START", "PostShipment", "NULL", "LAMBDA", message)
	utility.SendSQSLog(ctx, *message)

	device_id := message.Metadata["device_id"]

	body := map[string]interface{}{}
	err := json.Unmarshal(message.Body, &body)
	if err != nil {
		return err
	}

	// scrittura spedizione
	m := database.NewModel(models.NewShipment())
	m.Open()
	defer m.Close()

	var ok bool
	model := m.Model.(*models.Shipment)

	if model.OrderNo, ok = body["No"].(string); !ok {
		return errorFieldMissing("No")
	}
	if model.CustomerNo, ok = body["CustomerNo"].(string); !ok {
		return errorFieldMissing("CustomerNo")
	}
	if model.TotalWeight, ok = body["TotalWeight"].(float64); !ok {
		return errorFieldMissing("TotalWeight")
	}
	if model.Address, ok = body["Address"].(string); !ok {
		return errorFieldMissing("Address")
	}
	if model.City, ok = body["City"].(string); !ok {
		return errorFieldMissing("City")
	}
	if model.PostCode, ok = body["PostCode"].(string); !ok {
		return errorFieldMissing("PostCode")
	}
	if model.County, ok = body["County"].(string); !ok {
		return errorFieldMissing("County")
	}
	model.CustomerName, _ = body["CustomerName"].(string)
	model.Date = time.Now()
	model.VATRegistrationNo, _ = body["VATRegistrationNo"].(string)
	model.EMail, _ = body["EMail"].(string)
	model.PhoneNo, _ = body["PhoneNo"].(string)

	if !m.Create() {
		return m.GetLastError()
	}

	// cambia stato ordine
	var bodyBuff []byte
	bodyBuff, err = json.Marshal(map[string]string{
		"No": model.OrderNo,
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

	utility.BuildLogMetadata("END", "PostShipment", "NULL", "LAMBDA", message)
	utility.SendSQSLog(ctx, *message)
	return nil
}

func errorFieldMissing(field string) error {
	return errors.New("Field missing: " + field)
}

func main() {
	lambda.Start(func(ctx context.Context, snsEvent events.SNSEvent) error {
		return utility.HandlerSNSWithLogError(ctx, snsEvent, PostShipment)
	})

	// PostShipment(context.Background(), utility.Message{
	// 	Body: []byte(`{"No":"OR0001"}`),
	// 	Metadata: map[string]string{
	// 		"device_id": "eKtDSgRajJNmpcgKSQDTod:APA91bEE-t0pxZpNS5uG4jQEkV0I0P58fBBavNf9MldWxVp8xQJunv5UR7tReQ_seAWK1IxsdrinYANyqies47tmKpNStemyFTccwZiJJ8itsACJBQYlMVTw_rfG5ZO-_iIdgX5Y_QEZ",
	// 	},
	// })
}

package main

import (
	"context"
	"encoding/json"
	"errors"
	"thesis/financial-service/models"
	"thesis/lib/database"
	"thesis/lib/utility"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func PostInvoice(ctx context.Context, message *utility.Message) error {
	utility.BuildLogMetadata("START", "PostInvoice", "NULL", "LAMBDA", message)
	utility.SendSQSLog(ctx, *message)

	device_id := message.Metadata["device_id"]

	body := map[string]interface{}{}
	err := json.Unmarshal(message.Body, &body)
	if err != nil {
		return err
	}

	// scrittura spedizione
	m := database.NewModel(models.NewInvoiceHeader())
	m.Open()
	defer m.Close()

	var ok bool
	model := m.Model.(*models.InvoiceHeader)

	if model.OrderNo, ok = body["No"].(string); !ok {
		return errorFieldMissing("No")
	}
	if model.CustomerNo, ok = body["CustomerNo"].(string); !ok {
		return errorFieldMissing("CustomerNo")
	}
	if model.Amount, ok = body["Amount"].(float64); !ok {
		return errorFieldMissing("Amount")
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

	lines := body["Lines"].([]map[string]interface{})
	for i, line := range lines {
		m := database.NewModel(models.NewInvoiceLine())
		m.Open()

		modelLine := m.Model.(*models.InvoiceLine)
		modelLine.InvoiceNo = model.No
		modelLine.LineNo = i + 1
		modelLine.ItemNo = line["ItemNo"].(string)
		modelLine.ItemName = line["ItemName"].(string)
		modelLine.Quantity = line["Quantity"].(float64)
		modelLine.UnitPrice = line["UnitPrice"].(float64)
		modelLine.Amount = line["Amount"].(float64)

		if !m.Create() {
			m.Close()
			return m.GetLastError()
		}
		m.Close()
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
		"function":  "PostInvoice",
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

	utility.BuildLogMetadata("END", "PostInvoice", "NULL", "LAMBDA", message)
	utility.SendSQSLog(ctx, *message)
	return nil
}

func errorFieldMissing(field string) error {
	return errors.New("Field missing: " + field)
}

func main() {
	lambda.Start(func(ctx context.Context, snsEvent events.SNSEvent) error {
		return utility.HandlerSNSWithLogError(ctx, snsEvent, PostInvoice)
	})
}

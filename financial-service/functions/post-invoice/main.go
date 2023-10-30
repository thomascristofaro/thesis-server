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

	model := models.NewInvoiceHeader()
	if err = createInvoice(body, model); err != nil {
		return err
	}

	// cambia stato ordine
	var bodyBuff []byte
	bodyBuff, err = json.Marshal(map[string]string{
		"No":     model.OrderNo,
		"Status": "FATT",
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

func createInvoice(body map[string]interface{}, model *models.InvoiceHeader) error {
	m := database.NewModel(models.NewInvoiceHeader())
	m.Open()
	m.BeginTransaction()
	defer m.Close()

	var ok bool
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
	m.Model = model

	if !m.Create() {
		m.RollbackTransaction()
		return m.GetLastError()
	}

	mLine := database.NewModel(models.NewInvoiceLine())
	mLine.SetDB(m.GetDB())
	defer mLine.Close()

	lines := body["Lines"].([]interface{})
	for i, line := range lines {
		lineMap, ok := line.(map[string]interface{})
		if !ok {
			return errors.New("Line is not a map")
		}

		modelLine := models.NewInvoiceLine()
		modelLine.InvoiceID = model.ID
		modelLine.LineNo = i + 1
		if modelLine.ItemNo, ok = lineMap["ItemNo"].(string); !ok {
			mLine.RollbackTransaction()
			return errorFieldMissing("ItemNo")
		}
		if modelLine.Quantity, ok = lineMap["Quantity"].(float64); !ok {
			mLine.RollbackTransaction()
			return errorFieldMissing("Quantity")
		}
		if modelLine.UnitPrice, ok = lineMap["UnitPrice"].(float64); !ok {
			mLine.RollbackTransaction()
			return errorFieldMissing("UnitPrice")
		}
		if modelLine.Amount, ok = lineMap["Amount"].(float64); !ok {
			mLine.RollbackTransaction()
			return errorFieldMissing("Amount")
		}
		modelLine.ItemName, _ = lineMap["ItemName"].(string)
		mLine.Model = modelLine

		if !mLine.Create() {
			mLine.RollbackTransaction()
			return m.GetLastError()
		}
	}
	mLine.CommitTransaction()
	return nil
}

// TODO potrebbe diventare check field on map
func errorFieldMissing(field string) error {
	return errors.New("Field missing: " + field)
}

func main() {
	lambda.Start(func(ctx context.Context, snsEvent events.SNSEvent) error {
		return utility.HandlerSNSWithLogError(ctx, snsEvent, PostInvoice)
	})
}

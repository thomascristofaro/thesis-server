package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"thesis/lib/database"
	"thesis/lib/utility"
	"thesis/sales-service/models"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func ChangeStatusOrder(ctx context.Context, message *utility.Message) error {
	function := message.Metadata["function"]
	utility.BuildLogMetadata("START", "ChangeStatusOrder", "NULL", "LAMBDA", message)
	utility.SendSQSLog(ctx, *message)

	device_id := message.Metadata["device_id"]
	var status string

	switch function {
	case "PostShipment":
		status = "SPED"
	case "PostInvoice":
		status = "FATT"
	default:
		return errors.New(fmt.Sprintf("Event %s not handled", message.Metadata["event"]))
	}

	// cambio stato ordine
	body := map[string]interface{}{}
	err := json.Unmarshal(message.Body, &body)
	if err != nil {
		return err
	}
	fmt.Printf("Body: %v", body)
	no, ok := body["No"].(string)
	if !ok {
		return errors.New("Field missing: No")
	}

	m := database.NewModel(models.NewSalesOrderHeader())
	if !m.Open() {
		return m.GetLastError()
	}
	defer m.Close()
	m.SetFilter("No", database.EQUAL, no)
	if !m.Find() || !m.Next() {
		return m.GetLastError()
	}
	model := m.Model.(*models.SalesOrderHeader)
	model.Status = status
	if !m.Update() {
		return m.GetLastError()
	}

	// Notifica FCM
	utility.SendFirebaseNotification(ctx, device_id,
		"Change Status",
		fmt.Sprintf("Order %s status changed to %s", model.No, status),
		map[string]string{
			"action":  "update",
			"message": fmt.Sprintf("Order %s status changed to %s", model.No, status),
		},
	)

	// registra fattura
	if status == "SPED" {
		lines := []map[string]interface{}{}
		m := database.NewModel(models.NewSalesOrderLine())
		if !m.Open() {
			return m.GetLastError()
		}

		m.SetFilter("SalesOrderNo", database.EQUAL, model.No)
		if m.Find() {
			for m.Next() {
				line := m.Model.(*models.SalesOrderLine)
				lines = append(lines, map[string]interface{}{
					"ItemNo":    line.ItemNo,
					"ItemName":  line.ItemName,
					"Quantity":  line.Quantity,
					"UnitPrice": line.UnitPrice,
					"Amount":    line.Amount,
				})
			}
		}
		m.Close()

		body, err := json.Marshal(map[string]interface{}{
			"No":                model.No,
			"Amount":            model.Amount,
			"CustomerNo":        model.CustomerNo,
			"CustomerName":      model.CustomerName,
			"Address":           model.BillAddress,
			"City":              model.BillCity,
			"PostCode":          model.BillPostCode,
			"County":            model.BillCounty,
			"VATRegistrationNo": model.VATRegistrationNo,
			"EMail":             model.EMail,
			"PhoneNo":           model.PhoneNo,
			"Lines":             lines,
		})
		if err != nil {
			return err
		}
		message.Body = body

		// Send event to SNS
		err = utility.SendSNSMessage(
			context.Background(),
			"OnPostInvoiceTopicArn",
			*message)
		if err != nil {
			return err
		}
	}

	utility.BuildLogMetadata("END", "ChangeStatusOrder", "NULL", "LAMBDA", message)
	utility.SendSQSLog(ctx, *message)
	return nil
}

func main() {
	lambda.Start(func(ctx context.Context, sqsEvent events.SQSEvent) error {
		return utility.HandlerSQSWithLogError(ctx, sqsEvent, ChangeStatusOrder)
	})
}

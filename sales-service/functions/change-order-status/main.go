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
	utility.BuildLogMetadata("START", "ChangeStatusOrder", "NULL", "LAMBDA", message)
	utility.SendSQSLog(ctx, *message)

	device_id := message.Metadata["device_id"]
	var status string

	switch message.Metadata["function"] {
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

	m := database.NewModel(models.NewSalesOrderHeader())
	if !m.Open() {
		return m.GetLastError()
	}
	defer m.Close()
	m.SetFilter("No", database.EQUAL, body["No"].(string))
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
	// if status == "SPED" {
	// }

	utility.BuildLogMetadata("END", "ChangeStatusOrder", "NULL", "LAMBDA", message)
	utility.SendSQSLog(ctx, *message)
	return nil
}

func main() {
	lambda.Start(func(ctx context.Context, sqsEvent events.SQSEvent) error {
		return utility.HandlerSQSWithLogError(ctx, sqsEvent, ChangeStatusOrder)
	})
}

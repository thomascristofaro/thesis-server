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

func Handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	message, err := utility.ConvertSQSEventToMessage(sqsEvent)
	if err != nil {
		return err
	}

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
	err = json.Unmarshal(message.Body, &body)
	if err != nil {
		return err
	}

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
			"action": "update",
		},
	)

	// registra fattura
	// if status == "SPED" {
	// }

	return nil
}

func main() {
	lambda.Start(Handler)
}

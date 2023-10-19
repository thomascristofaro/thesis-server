package main

import (
	"context"
	"encoding/json"
	"fmt"
	"thesis/lib/database"
	"thesis/lib/utility"
	"thesis/utility-service/models"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func LogMessage(ctx context.Context, message utility.Message) error {
	m := database.NewModel(models.NewLog())

	status := message.Metadata["status"]
	function := message.Metadata["function"]
	event := message.Metadata["event"]
	service := message.Metadata["service"]
	device_id := message.Metadata["device_id"]
	uuid := message.Metadata["uuid"]
	delete(message.Metadata, "status")
	delete(message.Metadata, "function")
	delete(message.Metadata, "event")
	delete(message.Metadata, "service")
	delete(message.Metadata, "device_id")
	delete(message.Metadata, "uuid")

	b, _ := json.Marshal(message.Metadata)
	body := map[string]interface{}{}
	json.Unmarshal(message.Body, &body)
	message.Body, _ = json.Marshal(body["body"])

	m.Open()
	defer m.Close()
	model := m.Model.(*models.Log)
	model.Transaction = uuid
	model.Status = status
	model.Function = function
	model.Event = event
	model.Service = service
	model.Attributes = string(b)
	model.Body = string(message.Body)
	if !m.Create() {
		return m.GetLastError()
	}

	if status == "ERROR" && device_id != "" {
		err := utility.SendFirebaseNotification(ctx, device_id,
			"Error Log", "Error Log",
			map[string]string{
				"message": fmt.Sprintf("ERROR: %s", string(message.Body)),
			},
		)
		if err != nil {
			fmt.Printf("Errore invio notifica: %s", err.Error())
		}
	}
	return nil
}

func main() {
	lambda.Start(func(ctx context.Context, sqsEvent events.SQSEvent) error {
		message, err := utility.ConvertSQSEventToMessage(sqsEvent)
		if err != nil {
			return err
		}
		err = LogMessage(ctx, message)
		return nil
	})
}

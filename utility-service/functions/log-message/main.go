package main

import (
	"context"
	"encoding/json"
	"thesis/lib/database"
	"thesis/lib/utility"
	"thesis/utility-service/models"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func LogMessage(message utility.Message, function string, event string) error {
	m := database.NewModel(models.NewLog())
	b, _ := json.Marshal(message.Metadata)

	m.Open()
	model := m.Model.(*models.Log)
	model.Function = function
	model.Event = event
	model.Attributes = string(b)
	model.Body = string(message.Body)
	if !m.Create() {
		m.Close()
		return m.GetLastError()
	}
	m.Close()
	return nil
}

func main() {
	lambda.Start(func(ctx context.Context, sqsEvent events.SQSEvent) error {
		message, err := utility.ConvertSQSEventToMessage(sqsEvent)
		if err != nil {
			return err
		}
		err = LogMessage(message, "log-message", "on-log-message")
		return nil
	})
}

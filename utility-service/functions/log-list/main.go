package main

import (
	"context"
	"thesis/lib/component"
	"thesis/lib/utility"
	"thesis/utility-service/models"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type LogList struct {
	component.BasePage
	Model *models.Log
}

func NewLogList() *LogList {
	c := &LogList{
		Model: models.NewLog(),
	}
	c.BasePage = *component.NewBasePage(c, "Log List", component.PageList, "")
	c.AddModel(c.Model)
	c.AddKey("ID")
	c.AddSection(component.NewSection(
		"repeater",
		component.Repeater,
		"Repeater",
		component.NewField("ID", "ID", &c.Model.ID),
		component.NewField("Function", "Function", &c.Model.Function),
		component.NewField("Event", "Event", &c.Model.Event),
		component.NewField("Attributes", "Attributes", &c.Model.Attributes),
		component.NewField("Body", "Body", &c.Model.Body),
	))
	return c
}

func main() {
	lambda.Start(
		func(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
			return utility.AmazonHandler(ctx, request, NewLogList())
		})
}

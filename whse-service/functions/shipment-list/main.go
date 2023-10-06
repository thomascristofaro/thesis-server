package main

import (
	"context"
	"thesis/lib/component"
	"thesis/lib/utility"
	"thesis/whse-service/models"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(
		func(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
			return utility.AmazonHandler(ctx, request, NewShipmentList())
		})
}

type ShipmentList struct {
	component.BasePage
	Model *models.ShipmentHeader
}

func NewShipmentList() *ShipmentList {
	c := &ShipmentList{
		Model: models.NewShipmentHeader(),
	}
	c.BasePage = *component.NewBasePage(c, "Shipment List", component.PageList, "ShipmentCard")
	c.AddModel(c.Model)
	c.AddKey("No")
	c.AddSection(component.NewSection(
		"repeater1",
		component.Repeater,
		"List",
		component.NewField("No", "No.", &c.Model.No),
		component.NewField("CustomerNo", "Customer No.", &c.Model.CustomerNo),
		component.NewField("CustomerName", "Customer Name", &c.Model.CustomerName),
		component.NewField("Amount", "Amount", &c.Model.Amount),
	))
	return c
}

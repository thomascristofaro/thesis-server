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
	Model *models.Shipment
}

func NewShipmentList() *ShipmentList {
	c := &ShipmentList{
		Model: models.NewShipment(),
	}
	c.BasePage = *component.NewBasePage(c, "Shipment List", component.PageList, "ShipmentCard")
	c.AddModel(c.Model)
	c.AddKey("ID")
	c.AddSection(component.NewSection(
		"repeater1",
		component.Repeater,
		"List",
		component.NewField("ID", "ID", &c.Model.ID),
		component.NewField("OrderNo", "Order No.", &c.Model.OrderNo),
		component.NewField("CustomerNo", "Customer No.", &c.Model.CustomerNo),
		component.NewField("CustomerName", "Customer Name", &c.Model.CustomerName),
		component.NewField("Weight", "Weight", &c.Model.Weight),
	))
	return c
}

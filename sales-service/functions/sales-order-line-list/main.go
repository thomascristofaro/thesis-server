package main

import (
	"context"
	"thesis/lib/component"
	"thesis/lib/utility"
	"thesis/sales-service/models"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(
		func(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
			return utility.AmazonHandler(ctx, request, NewSalesOrderLineList())
		})
}

type SalesOrderLineList struct {
	component.BasePage
	Model *models.SalesOrderLine
}

func NewSalesOrderLineList() *SalesOrderLineList {
	c := &SalesOrderLineList{
		Model: models.NewSalesOrderLine(),
	}
	c.BasePage = *component.NewBasePage(c, "Lines", component.PageList, "SalesOrderLineCard")
	c.AddModel(c.Model)
	c.AddKey("SalesOrderNo")
	c.AddKey("LineNo")
	c.AddSection(component.NewSection(
		"repeater1",
		component.Repeater,
		"List",
		component.NewField("SalesOrderNo", "Sales Order No.", &c.Model.SalesOrderNo),
		component.NewField("LineNo", "Line No.", &c.Model.LineNo),
		component.NewField("ItemNo", "Item No.", &c.Model.ItemNo),
		component.NewField("ItemName", "Item Name", &c.Model.ItemName),
		component.NewField("Quantity", "Quantity", &c.Model.Quantity),
		component.NewField("UnitPrice", "Unit Price", &c.Model.UnitPrice),
		component.NewField("Amount", "Amount", &c.Model.Amount),
	))
	return c
}

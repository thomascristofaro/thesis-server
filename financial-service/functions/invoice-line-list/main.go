package main

import (
	"context"
	"thesis/financial-service/models"
	"thesis/lib/component"
	"thesis/lib/utility"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(
		func(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
			return utility.AmazonHandler(ctx, request, NewInvoiceLineList())
		})
}

type InvoiceLineList struct {
	component.BasePage
	Model *models.InvoiceLine
}

func NewInvoiceLineList() *InvoiceLineList {
	c := &InvoiceLineList{
		Model: models.NewInvoiceLine(),
	}
	c.BasePage = *component.NewBasePage(c, "Lines", component.PageList, "InvoiceLineCard")
	c.AddModel(c.Model)
	c.AddKey("InvoiceID")
	c.AddKey("LineNo")
	c.AddSection(component.NewSection(
		"repeater1",
		component.Repeater,
		"List",
		component.NewField("InvoiceID", "Invoice ID", &c.Model.InvoiceID),
		component.NewField("LineNo", "Line No.", &c.Model.LineNo),
		component.NewField("ItemNo", "Item No.", &c.Model.ItemNo),
		component.NewField("ItemName", "Item Name", &c.Model.ItemName),
		component.NewField("Quantity", "Quantity", &c.Model.Quantity),
		component.NewField("UnitPrice", "Unit Price", &c.Model.UnitPrice),
		component.NewField("Amount", "Amount", &c.Model.Amount),
	))
	return c
}

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
			return utility.AmazonHandler(ctx, request, NewInvoiceList())
		})
}

type InvoiceList struct {
	component.BasePage
	Model *models.InvoiceHeader
}

func NewInvoiceList() *InvoiceList {
	c := &InvoiceList{
		Model: models.NewInvoiceHeader(),
	}
	c.BasePage = *component.NewBasePage(c, "Invoice List", component.PageList, "InvoiceCard")
	c.AddModel(c.Model)
	c.AddKey("ID")
	c.AddSection(component.NewSection(
		"repeater1",
		component.Repeater,
		"List",
		component.NewField("ID", "ID", &c.Model.ID),
		component.NewField("CustomerNo", "Customer No.", &c.Model.CustomerNo),
		component.NewField("CustomerName", "Customer Name", &c.Model.CustomerName),
		component.NewField("Amount", "Amount", &c.Model.Amount),
	))
	return c
}

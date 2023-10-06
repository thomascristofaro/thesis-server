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
			return utility.AmazonHandler(ctx, request, NewCustomerList())
		})
}

type CustomerList struct {
	component.BasePage
	Model *models.Customer
}

func NewCustomerList() *CustomerList {
	c := &CustomerList{
		Model: models.NewCustomer(),
	}
	c.BasePage = *component.NewBasePage(c, "Customer List", component.PageList, "CustomerCard")
	c.AddModel(c.Model)
	c.AddKey("No")
	c.AddSection(component.NewSection(
		"repeater1",
		component.Repeater,
		"List",
		component.NewField("No", "No.", &c.Model.No),
		component.NewField("Name", "Name", &c.Model.Name),
		component.NewField("Balance", "Balance", &c.Model.Balance),
		component.NewField("VATRegistrationNo", "VAT Registration No.", &c.Model.VATRegistrationNo),
		component.NewField("EMail", "E-Mail", &c.Model.EMail),
	))
	return c
}

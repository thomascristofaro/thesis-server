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
			return utility.AmazonHandler(ctx, request, NewSalesOrderCard())
		})
}

type SalesOrderCard struct {
	component.BasePage
	Model *models.SalesOrderHeader
}

func NewSalesOrderCard() *SalesOrderCard {
	c := &SalesOrderCard{
		Model: models.NewSalesOrderHeader(),
	}
	c.BasePage = *component.NewBasePage(c, "Sales Order Card", component.PageCard, "")
	c.AddModel(c.Model)
	c.AddKey("No")
	c.AddButton(component.NewButton("post", "Registra", 0, "", nil))
	c.AddSection(component.NewSection(
		"Group1",
		component.Group,
		"General",
		component.NewField("No", "No.", &c.Model.No),
		component.NewField("Status", "Status", &c.Model.Status),
		component.NewField("CustomerNo", "Customer No.", &c.Model.CustomerNo),
		component.NewField("CustomerName", "Customer Name", &c.Model.CustomerName),
		component.NewField("VATRegistrationNo", "VAT Registration No.", &c.Model.VATRegistrationNo),
		// component.NewField("Date", "Date", &c.Model.Date),
		component.NewField("Amount", "Amount", &c.Model.Amount),
	))
	c.AddSection(component.NewSection(
		"Group2",
		component.Group,
		"Shipping & Invoice",
		component.NewField("Address", "Address", &c.Model.Address),
		component.NewField("City", "City", &c.Model.City),
		component.NewField("Post Code", "PostCode", &c.Model.PostCode),
		component.NewField("County", "County", &c.Model.County),
		component.NewField("EMail", "E-Mail", &c.Model.EMail),
		component.NewField("PhoneNo", "Phone No.", &c.Model.PhoneNo),
	))
	c.AddSection(component.NewSubPageSection(
		"Sub1",
		"Subpage 1",
		map[string]interface{}{
			"page_id": "SalesOrderLineList",
			"height":  400,
			"filters": []interface{}{
				map[string]interface{}{
					"id":    "SalesOrderNo",
					"value": 0,
					"field": "No",
				},
			},
		},
	))
	return c
}

package pages

import (
	"thesis/lib/component"
	"thesis/sales-service/models"
)

type SalesOrderList struct {
	component.BasePage
	Model *models.SalesOrderHeader
}

func NewSalesOrderList() *SalesOrderList {
	c := &SalesOrderList{
		Model: models.NewSalesOrderHeader(),
	}
	c.BasePage = *component.NewBasePage(c, "Sales Order List", component.PageList, "SalesOrderCard")
	c.AddModel(c.Model)
	c.AddKey("No")
	c.AddSection(component.NewSection(
		"repeater1",
		component.Repeater,
		"List",
		component.NewField("No", "No.", &c.Model.No),
		component.NewField("Status", "Status", &c.Model.Status),
		component.NewField("CustomerNo", "Customer No.", &c.Model.CustomerNo),
		component.NewField("CustomerName", "Customer Name", &c.Model.CustomerName),
		component.NewField("Amount", "Amount", &c.Model.Amount),
	))
	return c
}

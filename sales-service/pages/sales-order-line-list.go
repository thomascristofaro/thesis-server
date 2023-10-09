package pages

import (
	"thesis/lib/component"
	"thesis/sales-service/models"
)

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

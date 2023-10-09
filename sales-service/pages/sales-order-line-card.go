package pages

import (
	"thesis/lib/component"
	"thesis/sales-service/models"
)

type SalesOrderLineCard struct {
	component.BasePage
	Model *models.SalesOrderLine
}

func NewSalesOrderLineCard() *SalesOrderLineCard {
	c := &SalesOrderLineCard{
		Model: models.NewSalesOrderLine(),
	}
	c.BasePage = *component.NewBasePage(c, "Sales Order Line Card", component.PageCard, "")
	c.AddModel(c.Model)
	c.AddKey("No")
	c.AddSection(component.NewSection(
		"Group1",
		component.Group,
		"General",
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

package pages

import (
	"thesis/lib/component"
	"thesis/sales-service/models"
)

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

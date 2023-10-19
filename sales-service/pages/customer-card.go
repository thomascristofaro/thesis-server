package pages

import (
	"thesis/lib/component"
	"thesis/sales-service/models"
)

type CustomerCard struct {
	component.BasePage
	Model *models.Customer
}

func NewCustomerCard() *CustomerCard {
	c := &CustomerCard{
		Model: models.NewCustomer(),
	}
	c.BasePage = *component.NewBasePage(c, "Customer Card", component.PageCard, "")
	c.AddModel(c.Model)
	c.AddKey("No")
	c.AddSection(component.NewSection(
		"Group1",
		component.Group,
		"General",
		component.NewField("No", "No.", &c.Model.No),
		component.NewField("Name", "Name", &c.Model.Name),
		component.NewField("VATRegistrationNo", "VAT Registration No.", &c.Model.VATRegistrationNo),
		component.NewField("WebSite", "Web Site", &c.Model.WebSite),
		component.NewField("Balance", "Balance", &c.Model.Balance),
	))
	c.AddSection(component.NewSection(
		"Group2",
		component.Group,
		"Contacts",
		component.NewField("Address", "Address", &c.Model.Address),
		component.NewField("City", "City", &c.Model.City),
		component.NewField("Post Code", "PostCode", &c.Model.PostCode),
		component.NewField("County", "County", &c.Model.County),
		component.NewField("EMail", "E-Mail", &c.Model.EMail),
		component.NewField("PhoneNo", "Phone No.", &c.Model.PhoneNo),
	))
	return c
}

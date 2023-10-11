package pages

import (
	"context"
	"encoding/json"
	"errors"
	"thesis/lib/component"
	"thesis/lib/utility"
	"thesis/sales-service/models"

	"golang.org/x/exp/slices"
)

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
	c.AddButton(component.NewButton("post", "Registra", 0, "", PostSalesOrder))
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

func (p *SalesOrderCard) Button(queryParams map[string][]string) ([]byte, error) {
	var idSlices []string
	var ok bool
	if idSlices, ok = queryParams["button_id"]; !ok {
		return nil, errors.New("button_id not found in query parameters")
	}
	id := idSlices[0]
	bIdx := slices.IndexFunc(p.Buttons, func(b component.Button) bool { return b.Id == id })
	b := p.Buttons[bIdx]
	return b.Function(p, queryParams)
}

func PostSalesOrder(page component.Page, queryParams map[string][]string) ([]byte, error) {
	p := page.(*SalesOrderCard)

	// Handle Button
	device_id := queryParams["device_id"][0]
	delete(queryParams, "button_id")
	delete(queryParams, "device_id")

	// Search record
	if !p.ModelCtrl.Open() {
		return nil, p.ModelCtrl.GetLastError()
	}

	p.ModelCtrl.SetFilters(queryParams)
	if p.ModelCtrl.Find() {
		if !p.ModelCtrl.Next() {
			return nil, errors.New("No record found")
		}
	}
	p.ModelCtrl.Close()

	// Build message
	body, err := json.Marshal(map[string]interface{}{
		"No": p.Model.No,
	})
	if err != nil {
		return nil, err
	}
	metadata := map[string]string{
		"device_id": device_id,
		"function":  "PostSalesOrder",
		"service":   "SNS",
	}

	// registra spedizione
	if p.Model.Status == "" || p.Model.Status == "INIT" {
		metadata["event"] = "OnPostShipment"

		// TODO è da costruire il body per la spedizione

		// Send event to SNS
		ctx := context.Background()
		err = utility.SendSNSMessage(ctx,
			"OnPostShipmentTopicArn",
			utility.Message{
				Body:     body,
				Metadata: metadata,
			})
		if err != nil {
			return nil, err
		}
		return json.Marshal("In registrazione spedizione")
	}

	// registra fattura
	if p.Model.Status == "SPED" {
		metadata["event"] = "OnPostInvoice"

		// TODO è da costruire il body per la fattura

		// Send event to SNS
		ctx := context.Background()
		err = utility.SendSNSMessage(ctx,
			"OnPostInvoiceTopicArn",
			utility.Message{
				Body:     body,
				Metadata: metadata,
			})
		if err != nil {
			return nil, err
		}
		return json.Marshal("In registrazione fattura")
	}
	return json.Marshal("OK")
}

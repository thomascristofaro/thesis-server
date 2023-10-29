package pages

import (
	"context"
	"encoding/json"
	"errors"
	"thesis/lib/component"
	"thesis/lib/database"
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
		// component.NewField("Date", "Date", &c.Model.Date),
		component.NewField("Amount", "Amount", &c.Model.Amount),
	))
	c.AddSection(component.NewSection(
		"Group2",
		component.Group,
		"Customer Data",
		component.NewField("CustomerNo", "Customer No.", &c.Model.CustomerNo),
		component.NewField("CustomerName", "Customer Name", &c.Model.CustomerName),
		component.NewField("VATRegistrationNo", "VAT Registration No.", &c.Model.VATRegistrationNo),
		component.NewField("EMail", "E-Mail", &c.Model.EMail),
		component.NewField("PhoneNo", "Phone No.", &c.Model.PhoneNo),
	))
	c.AddSection(component.NewSection(
		"Group3",
		component.Group,
		"Shipping",
		component.NewField("ShipAddress", "Address", &c.Model.ShipAddress),
		component.NewField("ShipCity", "City", &c.Model.ShipCity),
		component.NewField("ShipPostCode", "Post Code", &c.Model.ShipPostCode),
		component.NewField("ShipCounty", "County", &c.Model.ShipCounty),
	))
	c.AddSection(component.NewSection(
		"Group4",
		component.Group,
		"Invoice",
		component.NewField("BillAddress", "Address", &c.Model.BillAddress),
		component.NewField("BillCity", "City", &c.Model.BillCity),
		component.NewField("BillPostCode", "Post Code", &c.Model.BillPostCode),
		component.NewField("BillCounty", "County", &c.Model.BillCounty),
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
	response := "OK"

	message := utility.Message{
		Body:     []byte("NULL"),
		Metadata: map[string]string{},
	}
	utility.BuildLogMetadata("START", "PostSalesOrder", "NULL", "LAMBDA", &message)
	utility.SendSQSLog(context.Background(), message)

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

	// Build general metadata
	message.Metadata["device_id"] = device_id

	// registra spedizione
	if p.Model.Status == "" || p.Model.Status == "INIT" {
		body, err := json.Marshal(map[string]interface{}{
			"No":                p.Model.No,
			"Weight":            p.Model.Weight,
			"CustomerNo":        p.Model.CustomerNo,
			"CustomerName":      p.Model.CustomerName,
			"Address":           p.Model.ShipAddress,
			"City":              p.Model.ShipCity,
			"PostCode":          p.Model.ShipPostCode,
			"County":            p.Model.ShipCounty,
			"VATRegistrationNo": p.Model.VATRegistrationNo,
			"EMail":             p.Model.EMail,
			"PhoneNo":           p.Model.PhoneNo,
		})
		if err != nil {
			return nil, err
		}
		message.Body = body

		// Send event to SNS
		err = utility.SendSNSMessage(
			context.Background(),
			"OnPostShipmentTopicArn",
			message)
		if err != nil {
			return nil, err
		}
		response = "In registrazione spedizione"
	}

	// registra fattura
	if p.Model.Status == "SPED" {
		lines := []map[string]interface{}{}
		m := database.NewModel(models.NewSalesOrderLine())
		if !m.Open() {
			return nil, m.GetLastError()
		}

		m.SetFilter("SalesOrderNo", database.EQUAL, p.Model.No)
		if m.Find() {
			for m.Next() {
				line := m.Model.(*models.SalesOrderLine)
				lines = append(lines, map[string]interface{}{
					"ItemNo":    line.ItemNo,
					"ItemName":  line.ItemName,
					"Quantity":  line.Quantity,
					"UnitPrice": line.UnitPrice,
					"Amount":    line.Amount,
				})
			}
		}
		m.Close()

		body, err := json.Marshal(map[string]interface{}{
			"No":                p.Model.No,
			"Amount":            p.Model.Amount,
			"CustomerNo":        p.Model.CustomerNo,
			"CustomerName":      p.Model.CustomerName,
			"Address":           p.Model.BillAddress,
			"City":              p.Model.BillCity,
			"PostCode":          p.Model.BillPostCode,
			"County":            p.Model.BillCounty,
			"VATRegistrationNo": p.Model.VATRegistrationNo,
			"EMail":             p.Model.EMail,
			"PhoneNo":           p.Model.PhoneNo,
			"Lines":             lines,
		})
		if err != nil {
			return nil, err
		}
		message.Body = body

		// Send event to SNS
		err = utility.SendSNSMessage(
			context.Background(),
			"OnPostInvoiceTopicArn",
			message)
		if err != nil {
			return nil, err
		}
		response = "In registrazione fattura"
	}

	utility.BuildLogMetadata("END", "PostSalesOrder", "NULL", "LAMBDA", &message)
	utility.SendSQSLog(context.Background(), message)
	return json.Marshal(response)
}

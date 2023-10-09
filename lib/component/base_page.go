package component

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"thesis/lib/database"

	"golang.org/x/exp/slices"
)

type BasePage struct {
	Id         string             `json:"id"`
	Type       PageType           `json:"type"`
	Caption    string             `json:"caption"`
	CardPageId string             `json:"card_page_id"`
	ModelCtrl  database.ModelCtrl `json:"-"`
	Keys       []string           `json:"key"`
	Buttons    []Button           `json:"buttons"`
	Area       []Section          `json:"area"`
}

func NewBasePage(pageImpl interface{}, caption string, pagetype PageType, CardPageId string) *BasePage {
	name := strings.ToLower(reflect.ValueOf(pageImpl).Elem().Type().Name())
	return &BasePage{
		Id:         name,
		Caption:    caption,
		Type:       pagetype,
		CardPageId: CardPageId,
		Buttons:    make([]Button, 0),
		Area:       make([]Section, 0),
		Keys:       make([]string, 0),
	}
}

func (p *BasePage) AddModel(model database.ModelDB) {
	p.ModelCtrl = database.NewModel(model)
}

func (p *BasePage) AddButton(b *Button) {
	p.Buttons = append(p.Buttons, *b)
}

func (p *BasePage) AddSection(s *Section) {
	p.Area = append(p.Area, *s)
}

func (p *BasePage) AddKey(s string) {
	p.Keys = append(p.Keys, s)
}

func (p *BasePage) GetId() string {
	return p.Id
}

func (p *BasePage) GetSchema() ([]byte, error) {
	return json.Marshal(p)
}

func (p *BasePage) Get(filters map[string][]string) ([]byte, error) {
	var recordset map[string]interface{}
	var data []map[string]interface{}
	if !p.ModelCtrl.Open() {
		return nil, p.ModelCtrl.GetLastError()
	}
	// section=XXXX <- mi deve dire anche il nome della sezione
	// filters=PK%3DCUSTOMER%3BSK%3D00002
	// No=123456&Name=THOMAS&Address=via+test+55&City=&PhoneNo=&VATRegistrationNo=
	p.ModelCtrl.SetFilters(filters)
	if p.ModelCtrl.Find() {
		for p.ModelCtrl.Next() {
			record := GetAllSectionsFieldsValue(p.Area)
			data = append(data, record)
		}
	}
	p.ModelCtrl.Close()
	recordset = make(map[string]interface{})
	recordset["recordset"] = data
	return json.Marshal(recordset)
}

func (p *BasePage) Post(body []byte) ([]byte, error) {
	if (body == nil) || (len(body) == 0) {
		return nil, errors.New("body is empty")
	}
	err := json.Unmarshal(body, p.ModelCtrl.Model)
	if err != nil {
		return nil, err
	}
	if !p.ModelCtrl.Open() {
		return nil, p.ModelCtrl.GetLastError()
	}
	if !p.ModelCtrl.Create() {
		return nil, p.ModelCtrl.GetLastError()
	}
	p.ModelCtrl.Close()
	return json.Marshal(p.ModelCtrl.Model)
}

func (p *BasePage) Patch(body []byte) ([]byte, error) {
	if (body == nil) || (len(body) == 0) {
		return nil, errors.New("body is empty")
	}
	err := json.Unmarshal(body, p.ModelCtrl.Model)
	if err != nil {
		return nil, err
	}
	if !p.ModelCtrl.Open() {
		return nil, p.ModelCtrl.GetLastError()
	}
	if !p.ModelCtrl.Update() {
		return nil, p.ModelCtrl.GetLastError()
	}
	p.ModelCtrl.Close()
	return json.Marshal(p.ModelCtrl.Model)
}

func (p *BasePage) Delete(filters map[string][]string) ([]byte, error) {
	if !p.ModelCtrl.Open() {
		return nil, p.ModelCtrl.GetLastError()
	}
	p.ModelCtrl.SetFilters(filters)
	if !(p.ModelCtrl.Find() && p.ModelCtrl.Next() && p.ModelCtrl.Delete()) {
		return nil, p.ModelCtrl.GetLastError()
	}
	p.ModelCtrl.Close()
	return json.Marshal("OK")
}

func (p *BasePage) Button(queryParams map[string][]string) ([]byte, error) {
	var idSlices []string
	var ok bool
	if idSlices, ok = queryParams["button_id"]; !ok {
		return nil, errors.New("button_id not found in query parameters")
	}
	id := idSlices[0]
	bIdx := slices.IndexFunc(p.Buttons, func(b Button) bool { return b.Id == id })
	b := p.Buttons[bIdx]
	err := b.Function(p, queryParams)
	if err != nil {
		return nil, err
	}
	return json.Marshal("OK")
}

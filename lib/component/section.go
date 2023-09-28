package component

import "reflect"

type Section struct {
	Id      string                 `json:"id"`
	Type    SectionType            `json:"type"`
	Caption string                 `json:"caption"`
	Fields  []Field                `json:"fields"`
	Options map[string]interface{} `json:"options"`
}

func NewSection(id string, t SectionType, caption string, fields ...Field) *Section {
	m := make([]Field, 0)
	for _, f := range fields {
		m = append(m, f)
	}
	return &Section{
		Id:      id,
		Type:    t,
		Caption: caption,
		Fields:  m,
		Options: map[string]interface{}{},
	}
}

func NewSubPageSection(id string, caption string, options map[string]interface{}) *Section {
	return &Section{
		Id:      id,
		Type:    SubPage,
		Caption: caption,
		Fields:  make([]Field, 0),
		Options: options,
	}
}

func GetSectionFieldsValue(s Section) map[string]interface{} {
	m := make(map[string]interface{})
	for _, f := range s.Fields {
		v := reflect.ValueOf(f.Value)
		if v.Kind() == reflect.Ptr {
			m[f.Id] = v.Elem().Interface()
		} else {
			m[f.Id] = f.Value
		}
	}
	return m
}

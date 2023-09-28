package component

import "fmt"

type Field struct {
	Id      string      `json:"id"`
	Caption string      `json:"caption"`
	Type    FieldType   `json:"type"`
	Value   interface{} `json:"-"`
}

func NewField(id string, caption string, value interface{}) Field {
	var fieldType FieldType
	switch v := value.(type) {
	case *string:
		fieldType = TextType
	case *int:
		fieldType = IntType
	case *float64:
		fieldType = DecimalType
	default:
		fmt.Printf("Tipo sconosciuto %T", v)
	}
	return Field{
		Id:      id,
		Caption: caption,
		Value:   value,
		Type:    fieldType,
	}
}

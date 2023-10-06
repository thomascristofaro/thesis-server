package models

import (
	"thesis/lib/database"
)

type InvoiceLine struct {
	InvoiceNo string `gorm:"primaryKey"`
	LineNo    int    `gorm:"primaryKey"`
	ItemNo    string
	ItemName  string
	Quantity  float64
	UnitPrice float64
	Amount    float64
}

func (c InvoiceLine) DBType() database.DBType {
	return database.SQL
}

func NewInvoiceLine() *InvoiceLine {
	return &InvoiceLine{}
}

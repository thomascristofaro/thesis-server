package models

import (
	"thesis/lib/database"
)

type InvoiceLine struct {
	InvoiceID uint `gorm:"column:invoiceid;primaryKey"`
	LineNo    int  `gorm:"column:lineno;primaryKey"`
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

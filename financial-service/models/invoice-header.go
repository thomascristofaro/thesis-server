package models

import (
	"thesis/lib/database"
	"time"
)

type InvoiceHeader struct {
	No           string `gorm:"primaryKey"`
	OrderNo      string
	CustomerNo   string
	CustomerName string
	Date         time.Time
	Amount       float64

	VATRegistrationNo string
	Address           string
	City              string
	PostCode          string
	County            string
	EMail             string
	PhoneNo           string
}

func (c InvoiceHeader) DBType() database.DBType {
	return database.SQL
}

func NewInvoiceHeader() *InvoiceHeader {
	return &InvoiceHeader{}
}

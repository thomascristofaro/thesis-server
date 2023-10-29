package models

import (
	"thesis/lib/database"
	"time"

	"gorm.io/gorm"
)

type InvoiceHeader struct {
	gorm.Model
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

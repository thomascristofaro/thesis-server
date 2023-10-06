package models

import (
	"thesis/lib/database"
	"time"
)

type SalesOrderHeader struct {
	No           string `gorm:"primaryKey"`
	CustomerNo   string
	CustomerName string
	Status       string
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

func (c SalesOrderHeader) DBType() database.DBType {
	return database.SQL
}

func NewSalesOrderHeader() *SalesOrderHeader {
	return &SalesOrderHeader{}
}

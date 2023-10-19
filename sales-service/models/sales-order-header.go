package models

import (
	"thesis/lib/database"
	"time"
)

type SalesOrderHeader struct {
	No                string `gorm:"primaryKey"`
	CustomerNo        string
	CustomerName      string
	Status            string
	Date              time.Time
	Amount            float64
	Weight            float64
	VATRegistrationNo string
	EMail             string
	PhoneNo           string

	ShipAddress  string
	ShipCity     string
	ShipPostCode string
	ShipCounty   string

	BillAddress  string
	BillCity     string
	BillPostCode string
	BillCounty   string
}

func (c SalesOrderHeader) DBType() database.DBType {
	return database.SQL
}

func NewSalesOrderHeader() *SalesOrderHeader {
	return &SalesOrderHeader{}
}

package models

import (
	"thesis/lib/database"
)

type Customer struct {
	No                string `gorm:"primaryKey"`
	Name              string
	VATRegistrationNo string
	Address           string
	City              string
	PostCode          string
	County            string
	EMail             string
	PhoneNo           string
	WebSite           string
	Balance           float64
	// vado a sommare quando nel change status order fatturo l'ordine
}

func (c Customer) DBType() database.DBType {
	return database.SQL
}

func NewCustomer() *Customer {
	return &Customer{}
}

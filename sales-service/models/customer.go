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
	Balance           float64 //Sum("Detailed Cust. Ledg. Entry".Amount WHERE("Customer No." = FIELD("No."),
}

func (c Customer) DBType() database.DBType {
	return database.SQL
}

func NewCustomer() *Customer {
	return &Customer{}
}

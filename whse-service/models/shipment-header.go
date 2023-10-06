package models

import (
	"thesis/lib/database"
	"time"
)

type ShipmentHeader struct {
	No           string `gorm:"primaryKey"`
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

func (c ShipmentHeader) DBType() database.DBType {
	return database.SQL
}

func NewShipmentHeader() *ShipmentHeader {
	return &ShipmentHeader{}
}

package models

import (
	"thesis/lib/database"
	"time"

	"gorm.io/gorm"
)

type Shipment struct {
	gorm.Model
	OrderNo      string
	CustomerNo   string
	CustomerName string
	Date         time.Time
	TotalWeight  float64

	VATRegistrationNo string
	Address           string
	City              string
	PostCode          string
	County            string
	EMail             string
	PhoneNo           string
}

func (c Shipment) DBType() database.DBType {
	return database.SQL
}

func NewShipment() *Shipment {
	return &Shipment{}
}

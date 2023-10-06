package models

import (
	"thesis/lib/database"
)

type ShipmentLine struct {
	ShipmentNo string `gorm:"primaryKey"`
	LineNo     int    `gorm:"primaryKey"`
	ItemNo     string
	ItemName   string
	Quantity   float64
	UnitPrice  float64
	Amount     float64
}

func (c ShipmentLine) DBType() database.DBType {
	return database.SQL
}

func NewShipmentLine() *ShipmentLine {
	return &ShipmentLine{}
}

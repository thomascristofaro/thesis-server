package models

import (
	"thesis/lib/database"
)

type SalesOrderLine struct {
	SalesOrderNo string `gorm:"primaryKey"`
	LineNo       int    `gorm:"primaryKey"`
	ItemNo       string
	ItemName     string
	Quantity     float64
	UnitPrice    float64
	Amount       float64
}

func (c SalesOrderLine) DBType() database.DBType {
	return database.SQL
}

func NewSalesOrderLine() *SalesOrderLine {
	return &SalesOrderLine{}
}

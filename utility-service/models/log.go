package models

import (
	"thesis/lib/database"

	"gorm.io/gorm"
)

type Log struct {
	gorm.Model
	Transaction string
	Status      string
	Function    string
	Event       string
	Service     string
	Attributes  string
	Body        string
}

func (c Log) DBType() database.DBType {
	return database.SQL
}

func NewLog() *Log {
	return &Log{}
}

package models

import (
	"thesis/lib/database"

	"gorm.io/gorm"
)

type Log struct {
	gorm.Model
	Function   string
	Event      string
	Attributes string
	Body       string
}

func (c Log) DBType() database.DBType {
	return database.SQL
}

func NewLog() *Log {
	return &Log{}
}

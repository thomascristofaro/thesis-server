package models

import "thesis/lib/database"

type Navigation struct {
	ID           uint `gorm:"primaryKey"`
	PageId       string
	Caption      string
	Tooltip      string
	Icon         int
	SelectedIcon int
	URL          string
}

func (c Navigation) DBType() database.DBType {
	return database.SQL
}

func NewNavigation() *Navigation {
	return &Navigation{}
}

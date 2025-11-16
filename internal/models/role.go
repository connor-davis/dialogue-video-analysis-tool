package models

import "github.com/lib/pq"

type Role struct {
	Base
	Name        string         `json:"name" gorm:"type:text;uniqueIndex;not null"`
	Description string         `json:"description" gorm:"type:text"`
	Permissions pq.StringArray `json:"permissions" gorm:"type:text[]"`
}

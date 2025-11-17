package models

import "github.com/google/uuid"

type Organization struct {
	Base
	Name    string    `json:"name" gorm:"type:text;not null"`
	Domain  string    `json:"domain" gorm:"type:text;not null"`
	OwnerId uuid.UUID `json:"ownerId" gorm:"type:uuid;not null"`
	Owner   User      `json:"owner" gorm:"foreignKey:OwnerId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Members []User    `json:"members" gorm:"many2many:organizations_members;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Roles   []Role    `json:"roles" gorm:"many2many:organizations_roles;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

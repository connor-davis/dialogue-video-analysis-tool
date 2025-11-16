package models

type User struct {
	Base
	Image       *string `json:"image" gorm:"type:text"`
	Name        string  `json:"name" gorm:"type:text;not null"`
	Username    string  `json:"username" gorm:"type:text;uniqueIndex;not null"`
	Bio         *string `json:"bio" gorm:"type:text"`
	MfaEnabled  bool    `json:"mfaEnabled" gorm:"type:boolean;default:false;not null"`
	MfaVerified bool    `json:"mfaVerified" gorm:"type:boolean;default:false;not null"`
	MfaSecret   []byte  `json:"-" gorm:"type:bytea"`
	Roles       []Role  `json:"-" gorm:"many2many:users_roles;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

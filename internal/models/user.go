package models

type UserType string

const (
	UserTypeSystemAdmin       UserType = "system_admin"
	UserTypeSystemUser        UserType = "system_user"
	UserTypeOrganizationOwner UserType = "organization_owner"
	UserTypeOrganizationUser  UserType = "organization_user"
)

type User struct {
	Base
	Name          string         `json:"name" gorm:"type:text;not null"`
	Email         string         `json:"email" gorm:"type:text;uniqueIndex;not null"`
	Password      []byte         `json:"-" gorm:"type:bytea"`
	Bio           *string        `json:"bio" gorm:"type:text"`
	MfaEnabled    bool           `json:"mfaEnabled" gorm:"type:boolean;default:false;not null"`
	MfaVerified   bool           `json:"mfaVerified" gorm:"type:boolean;default:false;not null"`
	MfaSecret     []byte         `json:"-" gorm:"type:bytea"`
	Type          UserType       `json:"type" gorm:"type:text;not null"`
	Roles         []Role         `json:"-" gorm:"many2many:users_roles;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Organizations []Organization `json:"-" gorm:"many2many:organizations_users;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

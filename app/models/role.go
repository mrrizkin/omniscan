package models

import (
	"time"

	"gorm.io/gorm"
)

type Role struct {
	ID              uint             `json:"id"                         gorm:"primary_key"`
	CreatedAt       *time.Time       `json:"created_at"`
	UpdatedAt       *time.Time       `json:"updated_at"`
	DeletedAt       gorm.DeletedAt   `json:"deleted_at"                 gorm:"index"`
	Slug            string           `json:"slug"                       gorm:"unique;not null;index"`
	Name            string           `json:"name"`
	RolePermissions []RolePermission `json:"role_permissions,omitempty" gorm:"foreignKey:RoleID"`
}

func (*Role) Seed(db *gorm.DB) {
	data := []Role{
		{Slug: "super_admin", Name: "Super Administrator"},
		{Slug: "user", Name: "User"},
	}

	for _, v := range data {
		db.FirstOrCreate(&Role{}, v)
	}
}

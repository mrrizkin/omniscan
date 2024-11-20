package models

import (
	"time"

	"gorm.io/gorm"
)

type RolePermission struct {
	ID           uint           `json:"id"                   gorm:"primary_key"`
	CreatedAt    *time.Time     `json:"created_at"`
	UpdatedAt    *time.Time     `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at"           gorm:"index"`
	RoleID       uint           `json:"role_id"`
	PermissionID uint           `json:"permission_id"`
	Role         *Role          `json:"role,omitempty"       gorm:"foreignKey:RoleID;references:ID"`
	Permission   *Permission    `json:"permission,omitempty" gorm:"foreignKey:PermissionID;references:ID"`
}

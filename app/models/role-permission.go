package models

import (
	"time"

	"github.com/mrrizkin/omniscan/system/utils"
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

func (*RolePermission) Seed(db *gorm.DB) {
	var roles []Role
	var permissions []Permission

	db.Find(&roles)
	db.Find(&permissions)

	if len(roles) == 0 || len(permissions) == 0 {
		return
	}

	for _, role := range roles {
		for _, permission := range permissions {
			if role.Slug == "super_admin" {
				if utils.Contains(permission.Slug,
					[]string{
						"create_permission",
						"read_permission",
						"update_permission",
						"delete_permission",

						"create_role_permission",
						"read_role_permission",
						"update_role_permission",
						"delete_role_permission",

						"create_role",
						"read_role",
						"update_role",
						"delete_role",

						"create_user",
						"read_user",
						"update_user",
						"delete_user",
					}) {
					db.FirstOrCreate(
						&RolePermission{},
						RolePermission{RoleID: role.ID, PermissionID: permission.ID},
					)
				}
			}

			if role.Slug == "user" {
				if utils.Contains(permission.Slug,
					[]string{
						"read_permission",

						"read_role_permission",

						"read_role",

						"read_user",
					}) {
					db.FirstOrCreate(
						&RolePermission{},
						RolePermission{RoleID: role.ID, PermissionID: permission.ID},
					)
				}
			}
		}
	}
}

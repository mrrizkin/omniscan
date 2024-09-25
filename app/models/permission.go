package models

import (
	"time"

	"gorm.io/gorm"
)

type Permission struct {
	ID        uint           `json:"id"         gorm:"primary_key"`
	CreatedAt *time.Time     `json:"created_at"`
	UpdatedAt *time.Time     `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	Group     string         `json:"group"`
	Slug      string         `json:"slug"       gorm:"unique;not null;index"`
	Name      string         `json:"name"`
}

func (*Permission) Seed(db *gorm.DB) {
	data := []Permission{
		{Group: "Permission", Slug: "create_permission", Name: "Create Permission"},
		{Group: "Permission", Slug: "read_permission", Name: "Read Permission"},
		{Group: "Permission", Slug: "update_permission", Name: "Update Permission"},
		{Group: "Permission", Slug: "delete_permission", Name: "Delete Permission"},

		{Group: "Role Permission", Slug: "create_role_permission", Name: "Create Role Permission"},
		{Group: "Role Permission", Slug: "read_role_permission", Name: "Read Role Permission"},
		{Group: "Role Permission", Slug: "update_role_permission", Name: "Update Role Permission"},
		{Group: "Role Permission", Slug: "delete_role_permission", Name: "Delete Role Permission"},

		{Group: "Role", Slug: "create_role", Name: "Create Role"},
		{Group: "Role", Slug: "read_role", Name: "Read Role"},
		{Group: "Role", Slug: "update_role", Name: "Update Role"},
		{Group: "Role", Slug: "delete_role", Name: "Delete Role"},

		{Group: "User", Slug: "create_user", Name: "Create User"},
		{Group: "User", Slug: "read_user", Name: "Read User"},
		{Group: "User", Slug: "update_user", Name: "Update User"},
		{Group: "User", Slug: "delete_user", Name: "Delete User"},

		{Group: "Special", Slug: "all", Name: "All"},
		{Group: "Special", Slug: "all_create", Name: "All Create"},
		{Group: "Special", Slug: "all_read", Name: "All Read"},
		{Group: "Special", Slug: "all_update", Name: "All Update"},
		{Group: "Special", Slug: "all_delete", Name: "All Delete"},
	}

	for _, v := range data {
		db.FirstOrCreate(&Permission{}, v)
	}
}

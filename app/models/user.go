package models

import (
	"time"

	"gorm.io/gorm"

	"github.com/mrrizkin/omniscan/system/config"
	"github.com/mrrizkin/omniscan/system/utils"
	"github.com/mrrizkin/omniscan/third-party/hashing"
)

type User struct {
	ID        uint           `json:"id"         gorm:"primary_key"`
	CreatedAt *time.Time     `json:"created_at"`
	UpdatedAt *time.Time     `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	Username  *string        `json:"username"   gorm:"unique;not null;index"`
	Password  *string        `json:"password"`
	Name      *string        `json:"name"`
	Email     *string        `json:"email"`
	RoleID    *uint          `json:"role_id"`
	Role      Role           `json:"role"       gorm:"foreignKey:RoleID;references:ID"`
}

func (*User) Seed(
	config *config.Config,
	hashing hashing.Hashing,
	db *gorm.DB,
) {
	var adminRole Role
	db.Where("slug = ?", "super_admin").First(&adminRole)

	username := config.SUPER_ADMIN_USERNAME
	password := config.SUPER_ADMIN_PASSWORD
	hash, err := hashing.GenerateHash(password)
	if err != nil {
		panic(err)
	}

	email := config.SUPER_ADMIN_EMAIL
	name := config.SUPER_ADMIN_NAME

	user := User{
		Username: &username,
		Password: &hash,
		Email:    &email,
		Name:     &name,
		RoleID:   &adminRole.ID,
	}

	userExist := new(User)
	wb := utils.NewWhereBuilder()
	wb.And("username = ?", username)
	wb.And("email = ?", email)
	wb.And("role_id = ?", adminRole.ID)
	wb.And("deleted_at IS NULL")
	where, whereArgs := wb.Get()

	err = db.Where(where, whereArgs...).
		First(userExist).
		Error

	if err == nil {
		return
	}

	db.Create(&user)
}

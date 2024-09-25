package user

import (
	"github.com/mrrizkin/omniscan/app/models"
	"github.com/mrrizkin/omniscan/system/database"
	"github.com/mrrizkin/omniscan/system/stypes"
)

func NewRepo(db *database.Database) *Repo {
	return &Repo{db}
}

func (r *Repo) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *Repo) FindAll(
	pagination stypes.Pagination,
) ([]models.User, error) {
	users := make([]models.User, 0)
	err := r.db.
		Preload("Role").
		Offset((pagination.Page - 1) * pagination.PerPage).
		Limit(pagination.PerPage).
		Find(&users).Error
	return users, err
}

func (r *Repo) FindAllCount() (int64, error) {
	var count int64 = 0
	err := r.db.Model(&models.User{}).Count(&count).Error
	return count, err
}

func (r *Repo) FindByID(id uint) (*models.User, error) {
	user := new(models.User)
	err := r.db.
		Preload("Role").
		First(user, id).
		Error
	return user, err
}

func (r *Repo) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *Repo) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}

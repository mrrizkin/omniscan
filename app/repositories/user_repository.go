package repositories

import (
	"github.com/mrrizkin/omniscan/app/models"
	"github.com/mrrizkin/omniscan/app/providers/database"
)

type UserRepository struct {
	db *database.Database
}

func (r *UserRepository) Construct() interface{} {
	return func(db *database.Database) *UserRepository {
		return &UserRepository{db}
	}
}

func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) FindAll(
	page int,
	perPage int,
) ([]models.User, error) {
	users := make([]models.User, 0)
	err := r.db.
		Preload("Role").
		Offset((page - 1) * perPage).
		Limit(perPage).
		Find(&users).Error
	return users, err
}

func (r *UserRepository) FindAllCount() (int64, error) {
	var count int64 = 0
	err := r.db.Model(&models.User{}).Count(&count).Error
	return count, err
}

var getUserPermissionsQuery = `SELECT p.*
FROM users u
JOIN roles r ON u.role_id = r.id
JOIN role_permissions rp ON r.id = rp.role_id
JOIN permissions p ON rp.permission_id = p.id
WHERE u.id = ?
ORDER BY p.id;`

func (r *UserRepository) UserPermissions(id uint) ([]models.Permission, error) {
	permissions := make([]models.Permission, 0)
	err := r.db.Raw(getUserPermissionsQuery).Scan(&permissions).Error
	return permissions, err
}

func (r *UserRepository) FindByID(id uint) (*models.User, error) {
	user := new(models.User)
	err := r.db.
		Preload("Role").
		First(user, id).
		Error
	return user, err
}

func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}

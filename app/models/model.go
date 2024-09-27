package models

import (
	"gorm.io/gorm"

	"github.com/mrrizkin/omniscan/system/config"
	"github.com/mrrizkin/omniscan/third-party/hashing"
)

type Seed interface {
	Seed(db *gorm.DB)
}

type Model struct {
	hashing hashing.Hashing
	config  *config.Config
	models  []interface{}
	seeds   []Seed
}

func New(config *config.Config, hashing hashing.Hashing) *Model {
	return &Model{
		models: []interface{}{
			&Permission{},
			&RolePermission{},
			&Role{},
			&User{},
			&Mutasi{},
			&MutasiDetail{},
		},
		seeds: []Seed{
			&Permission{},
			&Role{},
			&RolePermission{},
		},
		hashing: hashing,
		config:  config,
	}
}

func (m *Model) Migrate(db *gorm.DB) error {
	return db.AutoMigrate(m.models...)
}

func (m *Model) Seeds(db *gorm.DB) error {
	for _, model := range m.seeds {
		model.Seed(db)
	}

	userModel := new(User)
	userModel.Seed(m.config, m.hashing, db)

	return nil
}

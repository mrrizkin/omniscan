package provider

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/mrrizkin/omniscan/config"
)

type Mysql struct {
	config *config.Database
}

func NewMysql(config *config.Database) *Mysql {
	return &Mysql{config: config}
}

func (m *Mysql) DSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		m.config.USERNAME,
		m.config.PASSWORD,
		m.config.HOST,
		m.config.PORT,
		m.config.NAME,
	)
}

func (m *Mysql) Connect(cfg *config.Database) (*gorm.DB, error) {
	return gorm.Open(mysql.Open(m.DSN()))
}

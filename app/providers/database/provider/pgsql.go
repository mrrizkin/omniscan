package provider

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/mrrizkin/omniscan/config"
)

type Postgres struct {
	config *config.Database
}

func NewPostgres(config *config.Database) *Postgres {
	return &Postgres{config: config}
}

func (p *Postgres) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
		p.config.HOST,
		p.config.PORT,
		p.config.USERNAME,
		p.config.NAME,
		p.config.PASSWORD,
		p.config.SSLMODE,
	)
}

func (p *Postgres) Connect(cfg *config.Database) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(p.DSN()))
}

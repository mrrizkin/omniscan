package database

import (
	"context"
	"fmt"

	"go.uber.org/fx"
	"gorm.io/gorm"

	"github.com/mrrizkin/omniscan/app/providers/database/provider"
	"github.com/mrrizkin/omniscan/app/providers/logger"
	"github.com/mrrizkin/omniscan/config"
)

type DatabaseDriver interface {
	DSN() string
	Connect(cfg *config.Database) (*gorm.DB, error)
}

type Database struct {
	*gorm.DB
}

func (*Database) Construct() interface{} {
	return func(
		lc fx.Lifecycle,
		cfg *config.Database,
		log *logger.Logger,
	) (*Database, error) {
		var driver DatabaseDriver
		switch cfg.DRIVER {
		case "mysql", "mariadb", "maria":
			driver = provider.NewMysql(cfg)
		case "pgsql", "postgres", "postgresql":
			driver = provider.NewPostgres(cfg)
		case "sqlite", "sqlite3", "file":
			driver = provider.NewSqlite(cfg, log)
		default:
			return nil, fmt.Errorf("unknown database driver: %s", cfg.DRIVER)

		}

		gormDB, err := driver.Connect(cfg)
		if err != nil {
			return nil, err
		}

		db := Database{DB: gormDB}
		lc.Append(fx.Hook{
			OnStop: func(ctx context.Context) error {
				return db.Stop()
			},
		})

		return &db, nil
	}
}

func (d *Database) Stop() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}

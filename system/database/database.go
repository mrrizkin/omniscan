package database

import (
	"fmt"

	_ "github.com/joho/godotenv/autoload"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/mrrizkin/omniscan/app/models"
	"github.com/mrrizkin/omniscan/system/config"
	"github.com/mrrizkin/omniscan/third-party/logger"
)

type Database struct {
	*gorm.DB

	config *config.Config
	model  *models.Model
	logger logger.Logger
}

func New(config *config.Config, model *models.Model, logger logger.Logger) (*Database, error) {
	var (
		db  *gorm.DB
		err error
	)

	logger.Info("Connecting to database")

	switch config.DB_DRIVER {
	case "pgsql":
		db, err = gorm.Open(
			postgres.Open(
				fmt.Sprintf(
					"host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
					config.DB_HOST,
					config.DB_PORT,
					config.DB_USERNAME,
					config.DB_NAME,
					config.DB_PASSWORD,
					config.DB_SSLMODE,
				),
			),
		)
		if err != nil {
			return nil, err
		}
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(config.DB_HOST))
		if err != nil {
			return nil, err
		}

		err = db.Exec("PRAGMA journal_mode = WAL;").Error
		if err != nil {
			logger.Error(err, "Failed to enable WAL journal mode")
		} else {
			logger.Info("Enabled WAL journal mode")
		}

		err = db.Exec("PRAGMA foreign_keys = ON;").Error
		if err != nil {
			logger.Error(err, "Failed to enable foreign keys")
		} else {
			logger.Info("Enabled foreign keys")
		}
	case "mysql":
		db, err = gorm.Open(
			mysql.Open(
				fmt.Sprintf(
					"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
					config.DB_USERNAME,
					config.DB_PASSWORD,
					config.DB_HOST,
					config.DB_PORT,
					config.DB_NAME,
				),
			),
		)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", config.DB_DRIVER)
	}

	return &Database{
		DB:     db,
		logger: logger,
		config: config,
		model:  model,
	}, nil
}

func (d *Database) Start() error {
	if d.config.ENV != "prod" && d.config.ENV != "production" {
		d.logger.Info("Migrating model")
		err := d.model.Migrate(d.DB)
		if err != nil {
			return err
		}

		d.logger.Info("Seeding model")
		err = d.model.Seeds(d.DB)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *Database) Stop() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}

package provider

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/storage/mysql/v2"
	"github.com/gofiber/storage/postgres/v3"
	"github.com/gofiber/storage/sqlite3/v2"

	"github.com/mrrizkin/omniscan/config"
)

type Database struct {
	config *config.Database
}

func NewDatabase(config *config.Database) *Database {
	return &Database{config: config}
}

func (d *Database) Setup() (fiber.Storage, error) {
	switch d.config.DRIVER {
	case "pgsql", "postgres", "postgresql":
		return createPostgresStorage(d.config)
	case "mysql", "mariadb", "maria":
		return createMysqlStorage(d.config)
	case "sqlite", "sqlite3", "file":
		return createSQLiteStorage(d.config)
	default:
		return nil, fmt.Errorf("unknown database driver: %s", d.config.DRIVER)
	}
}

func createPostgresStorage(cfg *config.Database) (fiber.Storage, error) {
	config := postgres.Config{
		Host:       cfg.HOST,
		Port:       cfg.PORT,
		Database:   cfg.NAME,
		Username:   cfg.USERNAME,
		Password:   cfg.PASSWORD,
		Table:      "sessions",
		SSLMode:    cfg.SSLMODE,
		Reset:      false,
		GCInterval: 10 * time.Second,
	}

	return postgres.New(config), nil
}

func createMysqlStorage(cfg *config.Database) (fiber.Storage, error) {
	config := mysql.Config{
		Host:       cfg.HOST,
		Port:       cfg.PORT,
		Database:   cfg.NAME,
		Username:   cfg.USERNAME,
		Password:   cfg.PASSWORD,
		Table:      "sessions",
		Reset:      false,
		GCInterval: 10 * time.Second,
	}

	return mysql.New(config), nil
}

func createSQLiteStorage(cfg *config.Database) (fiber.Storage, error) {
	config := sqlite3.Config{
		Database:   cfg.HOST,
		Table:      "sessions",
		Reset:      false,
		GCInterval: 10 * time.Second,
	}

	return sqlite3.New(config), nil
}

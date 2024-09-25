package session

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/memory/v2"
	"github.com/gofiber/storage/mysql/v2"
	"github.com/gofiber/storage/postgres/v3"
	"github.com/gofiber/storage/sqlite3/v2"

	"github.com/mrrizkin/omniscan/system/config"
)

type Session struct {
	*session.Store

	storage fiber.Storage
}

func New(config *config.Config) (*Session, error) {
	var (
		storage fiber.Storage
		err     error
	)

	switch config.SESSION_DRIVER {
	case "database":
		storage, err = createDatabaseStorage(config)
	case "file":
		storage, err = createFileStorage()
	case "redis": // coming soon
		storage, err = createMemoryStorage()
	case "valkey": // coming soon
		storage, err = createMemoryStorage()
	case "memory":
		storage, err = createMemoryStorage()
	default:
		storage, err = createMemoryStorage()
	}

	if err != nil {
		return nil, err
	}

	return &Session{
		Store: session.New(session.Config{
			Storage:    storage,
			Expiration: 24 * time.Hour,
			KeyLookup:  "cookie:finteligo_session_key",
		}),
		storage: storage,
	}, nil
}

func (s *Session) Stop() error {
	return s.storage.Close()
}

func createDatabaseStorage(c *config.Config) (fiber.Storage, error) {
	switch c.DB_DRIVER {
	case "pgsql":
		return createPostgresStorage(c)
	case "mysql":
		return createMysqlStorage(c)
	case "sqlite":
		return createSQLiteStorage(c)
	default:
		return createMemoryStorage()
	}
}

func createPostgresStorage(c *config.Config) (fiber.Storage, error) {
	config := postgres.Config{
		Host:       c.DB_HOST,
		Port:       c.DB_PORT,
		Database:   c.DB_NAME,
		Table:      "sessions",
		SSLMode:    c.DB_SSLMODE,
		Reset:      false,
		GCInterval: 10 * time.Second,
	}

	return postgres.New(config), nil
}

func createMysqlStorage(c *config.Config) (fiber.Storage, error) {
	config := mysql.Config{
		Host:       c.DB_HOST,
		Port:       c.DB_PORT,
		Database:   c.DB_NAME,
		Table:      "sessions",
		Reset:      false,
		GCInterval: 10 * time.Second,
	}

	return mysql.New(config), nil
}

func createSQLiteStorage(c *config.Config) (fiber.Storage, error) {
	config := sqlite3.Config{
		Database:   c.DB_HOST,
		Table:      "sessions",
		Reset:      false,
		GCInterval: 10 * time.Second,
	}

	return sqlite3.New(config), nil
}

func createFileStorage() (fiber.Storage, error) {
	config := sqlite3.Config{
		Database:   "./storage/sessions.db",
		Table:      "sessions",
		Reset:      false,
		GCInterval: 10 * time.Second,
	}

	return sqlite3.New(config), nil
}

func createMemoryStorage() (fiber.Storage, error) {
	config := memory.Config{
		GCInterval: 10 * time.Second,
	}

	return memory.New(config), nil
}

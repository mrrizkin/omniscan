package provider

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/storage/sqlite3/v2"
)

type File struct{}

func NewFile() *File {
	return &File{}
}

func (*File) Setup() (fiber.Storage, error) {
	config := sqlite3.Config{
		Database:   "./storage/sessions.db",
		Table:      "sessions",
		Reset:      false,
		GCInterval: 10 * time.Second,
	}

	return sqlite3.New(config), nil
}

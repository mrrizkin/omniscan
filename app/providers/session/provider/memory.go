package provider

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/storage/memory/v2"
	"github.com/mrrizkin/omniscan/config"
)

type Memory struct {
	config *config.Session
}

func NewMemory(config *config.Session) *Memory {
	return &Memory{config: config}
}

func (m *Memory) Setup() (fiber.Storage, error) {
	switch m.config.DRIVER {
	case "memory":
		return createMemoryStorage()
	case "redis", "valkey":
		return nil, fmt.Errorf("driver %s is not yet supported", m.config.DRIVER)
	default:
		return nil, fmt.Errorf("unknown database driver: %s", m.config.DRIVER)
	}
}

func createMemoryStorage() (fiber.Storage, error) {
	config := memory.Config{
		GCInterval: 10 * time.Second,
	}

	return memory.New(config), nil
}

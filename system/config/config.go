package config

import (
	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	APP_NAME string `env:"APP_NAME,required"`
	APP_KEY  string `env:"APP_KEY,required"`
	ENV      string `env:"ENV,required"`
	PORT     int    `env:"PORT,required"`
	PREFORK  bool   `env:"PREFORK,default=false"`

	LOG_LEVEL      string `env:"LOG_LEVEL,default=debug"`
	LOG_CONSOLE    bool   `env:"LOG_CONSOLE,default=true"`
	LOG_FILE       bool   `env:"LOG_FILE,default=true"`
	LOG_DIR        string `env:"LOG_DIR,default=./storage/log"`
	LOG_MAX_SIZE   int    `env:"LOG_MAX_SIZE,default=50"`
	LOG_MAX_AGE    int    `env:"LOG_MAX_AGE,default=7"`
	LOG_MAX_BACKUP int    `env:"LOG_MAX_BACKUP,default=20"`
	LOG_JSON       bool   `env:"LOG_JSON,default=true"`

	HASH_MEMORY      int `env:"HASH_MEMORY,default=64"`
	HASH_ITERATIONS  int `env:"HASH_ITERATIONS,default=10"`
	HASH_PARALLELISM int `env:"HASH_PARALLELISM,default=2"`
	HASH_SALT_LEN    int `env:"HASH_SALT_LEN,default=32"`
	HASH_KEY_LEN     int `env:"HASH_KEY_LEN,default=32"`

	SUPER_ADMIN_NAME     string `env:"SUPER_ADMIN_NAME,required"`
	SUPER_ADMIN_EMAIL    string `env:"SUPER_ADMIN_EMAIL,required"`
	SUPER_ADMIN_USERNAME string `env:"SUPER_ADMIN_USERNAME,required"`
	SUPER_ADMIN_PASSWORD string `env:"SUPER_ADMIN_PASSWORD,required"`

	DB_DRIVER   string `env:"DB_DRIVER,default=sqlite"`
	DB_HOST     string `env:"DB_HOST,default=./storage/db.sqlite"`
	DB_PORT     int    `env:"DB_PORT,default=5432"`
	DB_NAME     string `env:"DB_NAME,default=omniscan"`
	DB_USERNAME string `env:"DB_USERNAME,default=omniscan"`
	DB_PASSWORD string `env:"DB_PASSWORD,default=omniscan"`
	DB_SSLMODE  string `env:"DB_SSLMODE,default=disable"`

	SESSION_DRIVER string `env:"SESSION_DRIVER,default=file"`
}

func New() (*Config, error) {
	config := new(Config)
	err := load(config)
	return config, err
}

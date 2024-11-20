package config

type Database struct {
	DRIVER       string `env:"DB_DRIVER,default=sqlite"`
	HOST         string `env:"DB_HOST"`
	PORT         int    `env:"DB_PORT,default=5432"`
	NAME         string `env:"DB_NAME"`
	USERNAME     string `env:"DB_USERNAME,default=root"`
	PASSWORD     string `env:"DB_PASSWORD,default=root"`
	SSLMODE      string `env:"DB_SSLMODE,default=disable"`
	AUTO_MIGRATE bool   `env:"DB_AUTO_MIGRATE,default=true"`
}

func (*Database) Construct() interface{} {
	return func() (*Database, error) {
		var database Database
		err := load(&database)
		return &database, err
	}
}

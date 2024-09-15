package gormfx

type Config struct {
	DBDriver   string `envconfig:"DB_DRIVER" default:""`
	DBUser     string `envconfig:"DB_USER" default:""`
	DBPassword string `envconfig:"DB_PASSWORD" default:""`
	DBHost     string `envconfig:"DB_HOST" default:""`
	DBName     string `envconfig:"DB_NAME" default:""`
	DBPort     string `envconfig:"DB_PORT" default:""`
}

func parseConfig() *Config {
	return &Config{
		DBDriver:   "mysql",
		DBUser:     "root",
		DBPassword: "password",
		DBHost:     "localhost",
		DBName:     "dittodining",
		DBPort:     "3306",
	}
}

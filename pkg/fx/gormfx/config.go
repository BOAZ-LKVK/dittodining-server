package gormfx

type Config struct {
	DBDriver               string `envconfig:"DB_DRIVER" default:""`
	DBUser                 string `envconfig:"DB_USER" default:""`
	DBPassword             string `envconfig:"DB_PASSWORD" default:""`
	DBHost                 string `envconfig:"DB_HOST" default:""`
	DBName                 string `envconfig:"DB_NAME" default:""`
	DBPort                 string `envconfig:"DB_PORT" default:""`
	ConnMaxLifetimeSeconds int64  `envconfig:"DB_CONN_MAX_LIFETIME_SECONDS" default:"0"`
	MaxIdleConns           int    `envconfig:"DB_MAX_IDLE_CONNS" default:"0"`
	MaxOpenConns           int    `envconfig:"DB_MAX_OPEN_CONNS" default:"0"`
}

func parseConfig() *Config {
	return &Config{
		DBDriver:               "mysql",
		DBUser:                 "root",
		DBPassword:             "password",
		DBHost:                 "localhost",
		DBName:                 "dittodining",
		DBPort:                 "3306",
		ConnMaxLifetimeSeconds: 1800,
		MaxIdleConns:           10, // TODO: DB 스펙에 따라 적절치로 조정
		MaxOpenConns:           10,
	}
}

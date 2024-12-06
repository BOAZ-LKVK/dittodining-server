package gormfx

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type Config struct {
	DBDriver               string `envconfig:"DB_DRIVER" default:"mysql"`
	DBUser                 string `envconfig:"DB_USER" default:"root"`
	DBPassword             string `envconfig:"DB_PASSWORD" default:"password"`
	DBHost                 string `envconfig:"DB_HOST" default:"localhost"`
	DBName                 string `envconfig:"DB_NAME" default:"dittodining"`
	DBPort                 string `envconfig:"DB_PORT" default:"3306"`
	ConnMaxLifetimeSeconds int64  `envconfig:"DB_CONN_MAX_LIFETIME_SECONDS" default:"1800"`
	MaxIdleConns           int    `envconfig:"DB_MAX_IDLE_CONNS" default:"10"` // TODO: DB 스펙에 따라 적절치로 조정
	MaxOpenConns           int    `envconfig:"DB_MAX_OPEN_CONNS" default:"10"`
}

func parseConfig() (*Config, error) {
	var config Config
	if err := envconfig.Process("", &config); err != nil {
		return nil, errors.New("failed to parse config from environment variables")
	}

	return &config, nil
}

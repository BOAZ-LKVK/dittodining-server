package fiberfx

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type Config struct {
	Port                 int    `envconfig:"PORT" default:"8080"`
	CORSAllowOrigins     string `envconfig:"CORS_ALLOW_ORIGINS" default:"*"`
	CORSAllowMethods     string `envconfig:"CORS_ALLOW_METHODS" default:"GET,POST,HEAD,PUT,DELETE,PATCH"`
	CORSAllowHeaders     string `envconfig:"CORS_ALLOW_HEADERS" default:"*"`
	CORSAllowCredentials bool   `envconfig:"CORS_ALLOW_CREDENTIALS" default:"false"`
}

func parseConfig() (*Config, error) {
	var config Config

	if err := envconfig.Process("", &config); err != nil {
		return nil, errors.New("failed to parse config from environment variables")
	}

	return &config, nil
}

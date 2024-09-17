package fiberfx

type Config struct {
	Port                 int    `envconfig:"PORT" default:"8080"`
	CORSAllowOrigins     string `envconfig:"CORS_ALLOW_ORIGINS" default:"*"`
	CORSAllowMethods     string `envconfig:"CORS_ALLOW_METHODS" default:"GET,POST,HEAD,PUT,DELETE,PATCH"`
	CORSAllowHeaders     string `envconfig:"CORS_ALLOW_HEADERS" default:""`
	CORSAllowCredentials bool   `envconfig:"CORS_ALLOW_CREDENTIALS" default:"false"`
}

func parseConfig() *Config {
	return &Config{
		Port:                 8080,
		CORSAllowOrigins:     "*",
		CORSAllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
		CORSAllowHeaders:     "*",
		CORSAllowCredentials: true,
	}
}

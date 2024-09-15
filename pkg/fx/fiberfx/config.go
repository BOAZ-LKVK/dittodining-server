package fiberfx

type Config struct {
	Port int `envconfig:"PORT" default:"8080"`
}

func parseConfig() *Config {
	return &Config{
		Port: 8080,
	}
}

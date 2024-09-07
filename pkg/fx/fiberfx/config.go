package fiberfx

type Config struct {
	Port int `envconfig:"PORT" default:"3000"`
}

func parseConfig() *Config {
	return &Config{
		Port: 3000,
	}
}

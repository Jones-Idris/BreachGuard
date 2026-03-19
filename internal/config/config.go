package config

import "os"

type Config struct {
	APIKey string
}

func Load() Config {
	return Config{
		APIKey: os.Getenv("HIBP_API_KEY"),
	}
}

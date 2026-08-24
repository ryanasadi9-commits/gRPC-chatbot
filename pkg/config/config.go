package config

import "os"

type Config struct {
	AppPort     string
	PostgresDSN string
	RedisAddr   string
}

func LoadEnv() *Config {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	return &Config{
		AppPort:     port,
		PostgresDSN: os.Getenv("POSTGRES_DSN"),
		RedisAddr:   os.Getenv("REDIS_ADDR"),
	}
}

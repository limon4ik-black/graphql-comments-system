package config

import "os"

type Config struct {
	AppPort     string
	PostgresDSN string
	RedisAddr   string
}

func Load() *Config {
	return &Config{
		AppPort:     getEnv("APP_PORT", "8080"),
		PostgresDSN: getEnv("POSTGRES_DSN", "postgres://graphql:graphql@localhost:5432/graphql_comments?sslmode=disable"),
		RedisAddr:   getEnv("REDIS_ADDR", "localhost:6379"),
	}
}

func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}

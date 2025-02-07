package config

import (
	"os"
)

type DBConfig struct {
	POSTGRES_DB       string
	POSTGRES_PASSWORD string
	POSTGRES_USER     string
	POSTGRES_PORT     string
	POSTGRES_HOST     string
}

type Config struct {
	DataBase DBConfig
}

func New() *Config {
	return &Config{
		DataBase: DBConfig{
			POSTGRES_DB:       getEnv("POSTGRES_DB", "pingdb"),
			POSTGRES_PASSWORD: getEnv("POSTGRES_PASSWORD", "12345"),
			POSTGRES_USER:     getEnv("POSTGRES_USER", "postgres"),
			POSTGRES_PORT:     getEnv("POSTGRES_PORT", "5432"),
			POSTGRES_HOST:     getEnv("POSTGRES_HOST", "db"),
		},
	}
}

func getEnv(key string, defaultVal string) string {
	if value, isExist := os.LookupEnv(key); isExist {
		return value
	}

	return defaultVal
}

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

type RabbitMQConfig struct {
	RABBITMQ_DEFAULT_USER string
	RABBITMQ_DEFAULT_PASS string
	RABBITMQ_PORT         string
}

type FrontConfig struct {
	FRONT_PORT string
}

type Config struct {
	DataBase       DBConfig
	FrontConfig    FrontConfig
	RabbitMQConfig RabbitMQConfig
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
		FrontConfig: FrontConfig{
			FRONT_PORT: getEnv("FRONT_PORT", "5137"),
		},
		RabbitMQConfig: RabbitMQConfig{
			RABBITMQ_DEFAULT_USER: getEnv("RABBITMQ_DEFAULT_USER", "user"),
			RABBITMQ_DEFAULT_PASS: getEnv("RABBITMQ_DEFAULT_PASS", "12345"),
			RABBITMQ_PORT:         getEnv("RABBITMQ_PORT", "5672"),
		},
	}
}

func getEnv(key string, defaultVal string) string {
	if value, isExist := os.LookupEnv(key); isExist {
		return value
	}

	return defaultVal
}

package config

import (
	"os"
	"strconv"
)

type PConfig struct {
	PINGER_WAIT int64
}
type RabbitMQConfig struct {
	RABBITMQ_DEFAULT_USER string
	RABBITMQ_DEFAULT_PASS string
	RABBITMQ_PORT         string
}

type Config struct {
	PConfig
	RabbitMQConfig
}

func New() *Config {
	return &Config{
		PConfig{
			PINGER_WAIT: getEnvAsInt("PINGER_WAIT", 7),
		},
		RabbitMQConfig{
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

func getEnvAsInt(key string, defaultVal int) int64 {
	valueStr := getEnv(key, strconv.Itoa(defaultVal))
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return int64(defaultVal)
	}
	return int64(value)
}

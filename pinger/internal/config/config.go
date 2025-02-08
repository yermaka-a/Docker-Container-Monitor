package config

import (
	"os"
	"strconv"
)

type PConfig struct {
	PINGER_WAIT int64
}

type Config struct {
	PConfig
}

func New() *Config {
	return &Config{
		PConfig{
			PINGER_WAIT: getEnvAsInt("PINGER_WAIT", 7),
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

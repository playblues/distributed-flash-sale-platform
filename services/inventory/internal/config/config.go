package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	RedisAddr string
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	return &Config{
		Port:      getEnv("PORT", "3001"),
		RedisAddr: getEnv("REDIS_ADDR", "localhost:6379"),
	}, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

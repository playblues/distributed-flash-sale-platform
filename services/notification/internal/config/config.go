package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	KafkaBroker string
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	return &Config{
		Port:        getEnv("PORT", "3002"),
		KafkaBroker: getEnv("KAFKA_BROKER", "localhost:9092"),
	}, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

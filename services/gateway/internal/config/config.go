package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                 string
	RedisAddr            string
	JWTSecret            string
	InventoryServiceURL  string
	OrderServiceURL      string
	InventoryGRPCEndpoint string
	OrderGRPCEndpoint     string
	RateLimit            int64
	RateLimitWindow      time.Duration
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	rateLimit, _ := strconv.ParseInt(getEnv("RATE_LIMIT", "5"), 10, 64)
	rateLimitWindow, _ := time.ParseDuration(getEnv("RATE_LIMIT_WINDOW", "10s"))

	return &Config{
		Port:                  getEnv("PORT", "3000"),
		RedisAddr:             getEnv("REDIS_ADDR", "localhost:6379"),
		JWTSecret:             getEnv("JWT_SECRET", "default_secret_key"),
		InventoryServiceURL:   getEnv("INVENTORY_SERVICE_URL", "http://localhost:3001"),
		OrderServiceURL:       getEnv("ORDER_SERVICE_URL", "http://localhost:8081"),
		InventoryGRPCEndpoint: getEnv("INVENTORY_GRPC_ENDPOINT", "localhost:50051"),
		OrderGRPCEndpoint:     getEnv("ORDER_GRPC_ENDPOINT", "localhost:9091"),
		RateLimit:             rateLimit,
		RateLimitWindow:       rateLimitWindow,
	}, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

package main

import (
	"context"
	"fmt"
	"log"

	"gateway-service/internal/config"
	"gateway-service/internal/handler"
	"gateway-service/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/redis/go-redis/v9"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Ошибка загрзуки конфигурации")
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: "",
		DB:       0,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {

		log.Fatalf("Ошибка подключения к Redis: %v", err)
	}
	fmt.Println("Успешное подключение к Redis")

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization, X-Requested-With, X-Idempotency-Key",
	}))

	app.Use(middleware.SlidingWindowLimiter(rdb, cfg.RateLimit, cfg.RateLimitWindow))
	app.Use(middleware.IdempotencyMiddleware(rdb))

	h := handler.NewGatewayHandler(cfg.JWTSecret, cfg.InventoryServiceURL, cfg.OrderServiceURL)

	app.Post("/api/v1/login", h.Login)
	app.Post("api/v1/buy", middleware.JWTMiddleware([]byte(cfg.JWTSecret)), h.Buy)

	fmt.Printf("Gateway запущен на http://localhost:%s", cfg.Port)
	log.Fatal(app.Listen(":" + cfg.Port))
}

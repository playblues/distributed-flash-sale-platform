package main

import (
	"context"
	"fmt"
	"log"

	"inventory-service/internal/config"
	"inventory-service/internal/handler"
	"inventory-service/internal/kafka"
	"inventory-service/internal/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/redis/go-redis/v9"
)

func main() {
	cfg, err := config.LoadConfig()

	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {

		log.Fatalf("Ошибка подключения к Redis: %v", err)
	}
	fmt.Println("Успешное подключение к Redis")

	repo := repository.NewInventoryRepository(rdb)
	h := handler.NewInventoryHandler(repo)

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))

	app.Get("/api/v1/inventory/stock", h.GetStock)
	app.Post("/api/v1/inventory/init", h.InitStock)
	app.Post("/api/v1/inventory/reserve", h.Reserve)
	app.Post("/api/v1/inventory/rollback", h.Rollback)

	go kafka.StartPaymentConsumer(repo)

	fmt.Printf("Inventory Service запущен на http://localhost:%s", cfg.Port)
	log.Fatal(app.Listen(":" + cfg.Port))
}

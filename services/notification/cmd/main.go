package main

import (
	"fmt"
	"log"

	"notification-service/internal/config"
	"notification-service/internal/kafka"
	"notification-service/internal/ws"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	hub := ws.NewHub()
	wsHandler := ws.NewHandler(hub)

	go kafka.StartConsumer(cfg.KafkaBroker, hub)

	app := fiber.New()

	app.Use("/ws", wsHandler.Upgrade)
	app.Get("/ws/:userId", websocket.New(wsHandler.HandleConnection))

	fmt.Printf("Notification Service запущен на порту %s\n", cfg.Port)
	log.Fatal(app.Listen(":" + cfg.Port))
}

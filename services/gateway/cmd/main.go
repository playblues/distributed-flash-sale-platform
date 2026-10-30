package main

import (
	"context"
	"fmt"
	"log"

	"gateway-service/internal/config"
	"gateway-service/internal/handler"
	"gateway-service/internal/middleware"
	inventoryv1 "gateway-service/pkg/api/inventory/v1"
	orderv1 "gateway-service/pkg/api/order/v1"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
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

	invConn, err := grpc.Dial(
		cfg.InventoryGRPCEndpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Не удалось подключиться к Inventory gRPC: %v", err)
	}
	defer invConn.Close()
	inventoryClient := inventoryv1.NewInventoryServiceClient(invConn)
	fmt.Printf("Подключен Inventory gRPC клиент (%s)\n", cfg.InventoryGRPCEndpoint)

	orderConn, err := grpc.Dial(
		cfg.OrderGRPCEndpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Не удалось подключиться к Order gRPC: %v", err)
	}
	defer orderConn.Close()
	orderClient := orderv1.NewOrderServiceClient(orderConn)
	fmt.Printf("Подключен Order gRPC клиент (%s)\n", cfg.OrderGRPCEndpoint)

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization, X-Requested-With, X-Idempotency-Key",
	}))

	app.Use(middleware.SlidingWindowLimiter(rdb, cfg.RateLimit, cfg.RateLimitWindow))
	app.Use(middleware.IdempotencyMiddleware(rdb))

	h := handler.NewGatewayHandler(cfg.JWTSecret, inventoryClient, orderClient)

	app.Post("/api/v1/login", h.Login)
	app.Post("/api/v1/buy", middleware.JWTMiddleware([]byte(cfg.JWTSecret)), h.Buy)

	fmt.Printf("Gateway запущен на http://localhost:%s\n", cfg.Port)
	log.Fatal(app.Listen(":" + cfg.Port))
}

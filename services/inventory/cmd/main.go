package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"inventory-service/internal/config"
	grpchandler "inventory-service/internal/grpc"
	"inventory-service/internal/handler"
	"inventory-service/internal/kafka"
	"inventory-service/internal/repository"
	inventoryv1 "inventory-service/pkg/api/inventory/v1"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
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

	go func() {
		lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
		if err != nil {
			log.Fatalf("Ошибка запуска gRPC listener на порту %s: %v", cfg.GRPCPort, err)
		}
		grpcServer := grpc.NewServer()
		inventoryv1.RegisterInventoryServiceServer(grpcServer, grpchandler.NewInventoryGrpcServer(repo))
		fmt.Printf("Inventory gRPC Server запущен на порту :%s\n", cfg.GRPCPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Ошибка работы gRPC сервера: %v", err)
		}
	}()

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

	fmt.Printf("Inventory HTTP Service запущен на http://localhost:%s\n", cfg.Port)
	log.Fatal(app.Listen(":" + cfg.Port))
}

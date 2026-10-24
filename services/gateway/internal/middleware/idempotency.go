package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

func IdempotencyMiddleware(rdb *redis.Client) fiber.Handler {
	ctx := context.Background()

	return func(c *fiber.Ctx) error {
		if c.Method() != fiber.MethodPost {
			return c.Next()
		}

		key := c.Get("X-Idempotency-Key")

		if key == "" {
			return c.Next()
		}

		redisKey := fmt.Sprintf("idempotency:%s", key)

		set, err := rdb.SetNX(ctx, redisKey, "PROCESSING", 30*time.Second).Result()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Ошибка проверки идемпотентности",
			})
		}

		if !set {
			cachedResponse, err := rdb.Get(ctx, redisKey).Result()
			if err == nil {
				if cachedResponse == "PROCESSING" {
					return c.Status(fiber.StatusConflict).JSON(fiber.Map{
						"error": "Запрос с данным ключом идемпотентности уже обрабатывается",
					})
				}
				c.Set("Content-Type", "application/json")
				c.Set("X-Cache", "HIT-IDEMPOTENT")
				return c.Status(fiber.StatusCreated).SendString(cachedResponse)
			}
		}

		if err := c.Next(); err != nil {
			rdb.Del(ctx, redisKey)
			return err
		}

		statusCode := c.Response().StatusCode()
		if statusCode == fiber.StatusCreated || statusCode == fiber.StatusOK {
			body := c.Response().Body()
			rdb.Set(ctx, redisKey, body, 24*time.Hour)
		} else {
			rdb.Del(ctx, redisKey)
		}
		return nil
	}
}

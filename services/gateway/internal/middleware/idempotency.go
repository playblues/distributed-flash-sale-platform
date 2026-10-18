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

		cachedResponse, err := rdb.Get(ctx, redisKey).Result()

		if err == nil && cachedResponse != "" {
			c.Set("Content-Type", "application/json")
			c.Set("X-cache", "HIT-IDEMPOTENT")
			return c.Status(fiber.StatusCreated).SendString(cachedResponse)
		}

		if err := c.Next(); err != nil {
			return err
		}
		if c.Response().StatusCode() == fiber.StatusCreated ||
			c.Response().StatusCode() == fiber.StatusOK {
			body := c.Response().Body()
			rdb.Set(ctx, redisKey, body, 24*time.Hour)
		}
		return nil
	}
}

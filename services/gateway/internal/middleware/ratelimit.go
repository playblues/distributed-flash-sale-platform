package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

func SlidingWindowLimiter(rdb *redis.Client, limit int64, window time.Duration) fiber.Handler {

	ctx := context.Background()
	return func(c *fiber.Ctx) error {
		ip := c.IP()
		key := fmt.Sprintf("rate_limit:sliding:%s", ip)

		limit := int64(5)
		window := 10 * time.Second

		now := time.Now()
		nowMS := now.UnixMilli()

		windowStartMs := now.Add(-window).UnixMilli()
		uniqueMember := fmt.Sprintf("%d", now.UnixNano())

		redisPipeline := rdb.TxPipeline()

		redisPipeline.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", windowStartMs))

		redisPipeline.ZAdd(ctx, key, redis.Z{
			Score:  float64(nowMS),
			Member: uniqueMember,
		})

		countCmd := redisPipeline.ZCard(ctx, key)

		redisPipeline.Expire(ctx, key, window)

		_, err := redisPipeline.Exec(ctx)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal Redis Error",
			})
		}
		count := countCmd.Val()

		if count > limit {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": fmt.Sprintf("Слшком много запросов. Сделано : %d, Лимит: %d", count, limit),
			})
		}
		return c.Next()
	}
}

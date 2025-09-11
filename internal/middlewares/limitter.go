package middlewares

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"

	"github.com/bagasunix/transnovasi/pkg/env"
)

func HybridRateLimiter(redisClient *redis.Client, cfg *env.Cfg) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext() // pakai context Fiber
		ip := c.IP()
		now := time.Now().Unix()

		// ----- TOKEN BUCKET -----
		bucketKey := fmt.Sprintf("bucket:%s", ip)
		// isi ulang token (refill_rate token per detik)
		refillRate := float64(cfg.Server.RateLimiter.Limit) / cfg.Server.RateLimiter.Duration.Seconds()
		capacity := cfg.Server.RateLimiter.Limit

		vals, err := redisClient.HMGet(ctx, bucketKey, "tokens", "ts").Result()
		if err != nil && err != redis.Nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":    fiber.StatusInternalServerError,
				"error":   "server_error",
				"message": "cannot fetch token bucket from redis",
			})
		}

		tokens, _ := strconv.ParseFloat(fmt.Sprint(vals[0]), 64)
		lastTs, _ := strconv.ParseInt(fmt.Sprint(vals[1]), 10, 64)
		if lastTs == 0 {
			tokens = float64(capacity)
			lastTs = now
		}

		elapsed := now - lastTs
		tokens = math.Min(float64(capacity), tokens+float64(elapsed)*refillRate)
		lastTs = now

		if tokens < 1 {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"code":    fiber.StatusTooManyRequests,
				"error":   "rate_limit",
				"message": "token bucket empty, too many requests",
			})
		}
		tokens--

		pipe := redisClient.TxPipeline()
		// Simpan kembali
		pipe.HSet(ctx, bucketKey, "tokens", tokens, "ts", lastTs)
		pipe.Expire(ctx, bucketKey, cfg.Server.RateLimiter.Duration*2)
		if _, err := pipe.Exec(ctx); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":    fiber.StatusInternalServerError,
				"error":   "server_error",
				"message": "cannot update token bucket in redis",
			})
		}

		// ----- SLIDING WINDOW -----
		winSize := int64(cfg.Server.RateLimiter.Duration.Seconds())
		currWin := now / winSize
		prevWin := currWin - 1
		keyCurr := fmt.Sprintf("sw:%s:%d", ip, currWin)
		keyPrev := fmt.Sprintf("sw:%s:%d", ip, prevWin)

		currCount, err := redisClient.Incr(ctx, keyCurr).Result()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":    fiber.StatusInternalServerError,
				"error":   "server_error",
				"message": "cannot increment sliding window",
			})
		}
		if currCount == 1 {
			redisClient.Expire(ctx, keyCurr, cfg.Server.RateLimiter.Duration*2)
		}

		prevCount, err := redisClient.Get(ctx, keyPrev).Int64()
		if err != nil && err != redis.Nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":    fiber.StatusInternalServerError,
				"error":   "server_error",
				"message": "cannot fetch previous sliding window",
			})
		}

		elapsedWin := now % winSize
		est := float64(currCount) + float64(prevCount)*(1-float64(elapsedWin)/float64(winSize))

		if est > float64(cfg.Server.RateLimiter.Limit) {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"code":    fiber.StatusTooManyRequests,
				"error":   "rate_limit",
				"message": "sliding window exceeded, too many requests",
			})
		}

		// ----- ADAPTIVE -----
		// kalau request > 3x limit → penalti (turunkan refillRate sementara)
		if est > float64(cfg.Server.RateLimiter.Limit*3) {
			redisClient.Set(ctx, fmt.Sprintf("penalty:%s", ip), "1", 5*time.Minute)
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"code":    fiber.StatusTooManyRequests,
				"error":   "rate_limit",
				"message": "adaptive block triggered, request limit exceeded",
			})
		}

		return c.Next()
	}
}

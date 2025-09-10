package middlewares

import (
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/phuslu/log"

	"github.com/bagasunix/transnovasi/pkg/env"
	"github.com/bagasunix/transnovasi/pkg/jwt"
)

func NewAuthMiddleware(redisClient *redis.Client, log *log.Logger, cfg *env.Cfg) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		authHeader := ctx.Get("Authorization")
		if authHeader == "" {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing authorization header"})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid authorization header format"})
		}

		// Gunakan ValidateToken untuk memvalidasi JWT
		claims, err := jwt.ValidateToken(log, tokenString, cfg)
		if err != nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}

		var redisKey string
		if claims.User.ID != 0 {
			redisKey = "auth_user:token:" + strconv.Itoa(int(claims.User.ID))
			// Simpan claims ke context jika diperlukan untuk endpoint lain
			ctx.Locals("user", claims)
		} else {
			redisKey = "auth_customer:token:" + strconv.Itoa(int(claims.Customer.ID))
			ctx.Locals("customer", claims)
		}

		// Periksa apakah token ada di Redis
		exists, err := redisClient.Exists(ctx.Context(), redisKey).Result()
		if err != nil || exists == 0 {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "session expired, please login again"})
		}

		// Cek sisa waktu expired token di Redis
		ttl, err := redisClient.TTL(ctx.Context(), tokenString).Result()
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "error checking session"})
		}

		// Perpanjang masa aktif token jika kurang dari 1 menit sebelum kedaluwarsa
		if ttl > 0 && ttl < time.Minute {
			redisClient.Expire(ctx.Context(), tokenString, 15*time.Minute)
		}

		return ctx.Next()
	}
}

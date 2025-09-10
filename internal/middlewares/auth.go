package middlewares

import (
	"context"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/phuslu/log"

	"github.com/bagasunix/transnovasi/internal/dtos/responses"
	"github.com/bagasunix/transnovasi/pkg/env"
	"github.com/bagasunix/transnovasi/pkg/jwt"
)

type ctxKey string

const authClaimsKey ctxKey = "authorization_payload"

func NewAuthMiddleware(redisClient *redis.Client, log *log.Logger, cfg *env.Cfg) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		authHeader := ctx.Get("Authorization")
		if authHeader == "" {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Anda belum login", "error": "missing authorization header"})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Anda belum login", "error": "invalid authorization header format"})
		}

		// Gunakan ValidateToken untuk memvalidasi JWT
		claims, err := jwt.ValidateToken(log, tokenString, cfg)
		if err != nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Anda belum login", "error": err.Error()})
		}

		var redisKey string = "0"
		redisKey = "auth_user:token:" + claims.User.ID

		// Periksa apakah token ada di Redis
		exists, err := redisClient.Exists(ctx.Context(), redisKey).Result()
		if err != nil || exists == 0 {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Anda belum login", "error": "session expired, please login again"})
		}

		// Cek sisa waktu expired token di Redis
		ttl, err := redisClient.TTL(ctx.Context(), tokenString).Result()
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Anda belum login", "error": "error checking session"})
		}

		// Perpanjang masa aktif token jika kurang dari 1 menit sebelum kedaluwarsa
		if ttl > 0 && ttl < time.Minute {
			redisClient.Expire(ctx.Context(), tokenString, 15*time.Minute)
		}

		// 1. Set ke Fiber Locals (untuk akses langsung dari fiber.Ctx)
		ctx.Locals(string(authClaimsKey), claims.User)

		// 2. Set ke User Context (untuk akses dari context.Context di usecase)
		userCtx := context.WithValue(ctx.UserContext(), authClaimsKey, claims.User)
		ctx.SetUserContext(userCtx)

		return ctx.Next()
	}
}

// Helper function untuk mengambil claims dari fiber context
func GetAuthClaims(ctx *fiber.Ctx) *responses.UserResponse {
	claims := ctx.Locals(string(authClaimsKey))
	if claims == nil {
		return nil
	}

	authClaims, _ := claims.(*responses.UserResponse)
	return authClaims
}

// Helper function untuk mengambil claims dari context.Context (untuk usecase)
func GetAuthClaimsFromContext(ctx context.Context) *responses.UserResponse {
	claims := ctx.Value(string(authClaimsKey))
	if claims == nil {
		return nil
	}

	authClaims, _ := claims.(*responses.UserResponse)
	return authClaims
}

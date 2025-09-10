package middlewares

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/phuslu/log"

	"github.com/bagasunix/transnovasi/internal/domains"
	"github.com/bagasunix/transnovasi/internal/repositories"
	"github.com/bagasunix/transnovasi/pkg/jwt"
)

// Middleware generator
func LoggingMiddleware(logger *log.Logger, r repositories.Repositories) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// ambil request body (hati-hati kalau besar, bisa di-limit)
		reqBody := string(c.Body())

		// jalankan handler berikutnya
		err := c.Next()

		// ambil response body
		resBody := string(c.Response().Body())

		// ambil status code
		statusCode := strconv.Itoa(c.Response().StatusCode())

		// ambil user & customer ID dari locals (kalau ada)
		userID := "0"
		if c.Locals("user") != nil {
			userID = strconv.Itoa(int(c.Locals("user").(*jwt.Claims).User.ID))
		}

		// simpan ke database
		logEntry := domains.AuditLog{
			UserID:     userID,
			Method:     c.Method(),
			Endpoint:   c.OriginalURL(),
			StatusCode: statusCode,
			Request:    reqBody,
			Response:   resBody,
			UserAgent:  c.Get("User-Agent"),
			IPAddress:  c.IP(),
			CreatedAt:  time.Now(),
		}
		if saveErr := r.GetAuditLog().Create(c.Context(), &logEntry); saveErr != nil {
			logger.Error().Err(saveErr).Msg("Failed to save audit log")
		}

		// log ke console
		logger.Info().
			Str("method", c.Method()).
			Str("endpoint", c.OriginalURL()).
			Int("status", c.Response().StatusCode()).
			Str("user_agent", c.Get("User-Agent")).
			Str("ip_address", c.IP()).
			Dur("duration", time.Since(start)).
			Msg(resBody)

		return err
	}
}

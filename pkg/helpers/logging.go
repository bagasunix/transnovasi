package helpers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/phuslu/log"

	"github.com/bagasunix/transnovasi/internal/domains"
	"github.com/bagasunix/transnovasi/internal/repositories"
	"github.com/bagasunix/transnovasi/pkg/jwt"
)

func LoggingMiddleware(c *fiber.Ctx, log *log.Logger, r repositories.Repositories, reqBody, resBody, statusCOde string) error {
	start := time.Now()

	var userID string = "0"

	if c.Locals("user") != nil {
		userID = strconv.Itoa(int(c.Locals("user").(*jwt.Claims).User.ID))
	}
	// Simpan log ke database
	logEntry := domains.AuditLog{
		UserID:     userID,
		Method:     c.Method(),
		Endpoint:   c.OriginalURL(),
		StatusCode: statusCOde,
		Request:    reqBody,
		Response:   resBody,
		UserAgent:  fmt.Sprintf("%v", c.Get("User-Agent")), // Add user agent
		IPAddress:  c.IP(),                                 // Add IP address
		CreatedAt:  time.Now(),
	}

	if err := r.GetAuditLog().Create(c.Context(), &logEntry); err != nil {
		log.Error().Err(err).Msg("Failed to save audit log")
	}

	log.Info().
		Str("method", c.Method()).
		Str("endpoint", c.OriginalURL()).
		Int("status", c.Response().StatusCode()).
		Str("user_agent", fmt.Sprintf("%v", c.Get("User-Agent"))). // Log user agent
		Str("ip_address", c.IP()).                                 // Log IP address
		Dur("duration", time.Since(start)).
		Msg(resBody)

	return nil
}

package configs

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/recover"
	fiberSwagger "github.com/swaggo/fiber-swagger"

	_ "github.com/bagasunix/transnovasi/docs" // docs yang digenerate swag
	"github.com/bagasunix/transnovasi/pkg/env"
)

func InitFiber(ctx context.Context, cfg *env.Cfg) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName: cfg.App.Name,
	})
	app.Use(func(c *fiber.Ctx) error {
		c.Set("Content-Security-Policy", "default-src 'self'")
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		return c.Next()
	})
	app.Use(helmet.New())
	app.Use(recover.New())
	app.Use(favicon.New())

	// route swagger UI
	app.Get("swagger/*", fiberSwagger.WrapHandler)
	return app
}

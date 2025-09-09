package http

import (
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"

	"github.com/bagasunix/transnovasi/internal/controllers"
	"github.com/bagasunix/transnovasi/internal/delivery/http/handlers"
	"github.com/bagasunix/transnovasi/pkg/env"
)

type RouteConfig struct {
	App            *fiber.App
	AuthController *controllers.AuthController
	Cfg            *env.Cfg
	Redis          *redis.Client
	Logger         *log.Logger
}

func InitHttpHandler(f *RouteConfig) *fiber.App {
	return NewHttpHandler(*f)
}

func NewHttpHandler(r RouteConfig) *fiber.App {
	// Initialize middleware
	// Handlers
	handlers.MakeAuthHandler(r.AuthController, r.App.Group(r.Cfg.Server.Version+"/auth").(*fiber.Group))
	return r.App
}

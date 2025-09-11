package http

import (
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/phuslu/log"

	"github.com/bagasunix/transnovasi/internal/controllers"
	"github.com/bagasunix/transnovasi/internal/delivery/http/handlers"
	"github.com/bagasunix/transnovasi/internal/middlewares"
	"github.com/bagasunix/transnovasi/pkg/env"
)

type RouteConfig struct {
	App               *fiber.App
	AuthController    *controllers.AuthController
	CustController    *controllers.CustomerController
	VehicleController *controllers.VehicleController
	Cfg               *env.Cfg
	Redis             *redis.Client
	Logger            *log.Logger
}

func InitHttpHandler(f *RouteConfig) *fiber.App {
	return NewHttpHandler(*f)
}

func NewHttpHandler(r RouteConfig) *fiber.App {
	// Initialize middleware
	authMiddleware := middlewares.NewAuthMiddleware(r.Redis, r.Logger, r.Cfg)
	// Handlers
	handlers.MakeAuthHandler(r.AuthController, r.App.Group(r.Cfg.Server.Version+"/auth").(*fiber.Group), authMiddleware)
	handlers.MakeCustHandler(r.CustController, r.App.Group(r.Cfg.Server.Version+"/customer").(*fiber.Group), authMiddleware)
	handlers.MakeVehicleHandler(r.VehicleController, r.App.Group(r.Cfg.Server.Version+"/vehicle").(*fiber.Group), authMiddleware)
	return r.App
}

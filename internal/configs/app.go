package configs

import (
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/phuslu/log"
	"gorm.io/gorm"

	"github.com/bagasunix/transnovasi/internal/controllers"
	"github.com/bagasunix/transnovasi/internal/delivery/http"
	"github.com/bagasunix/transnovasi/internal/repositories"
	"github.com/bagasunix/transnovasi/internal/usecases"
	"github.com/bagasunix/transnovasi/pkg/env"
)

type setupApp struct {
	DB    *gorm.DB
	App   *fiber.App
	Log   *log.Logger
	Cfg   *env.Cfg
	Redis *redis.Client
}

func SetupApp(app *setupApp) *http.RouteConfig {
	app.Log.Info().Msg("Setting up application...")
	// setup repositories
	repositories := repositories.New(app.Log, app.DB)

	// setup use cases
	authUsecase := usecases.NewAuthUsecase(app.Log, app.DB, app.Cfg, repositories, app.Redis)
	// setup controller
	authContoller := controllers.NewAuthController(app.Log, repositories, authUsecase)

	return &http.RouteConfig{
		App:            app.App,
		AuthController: authContoller,
		Cfg:            app.Cfg,
		Redis:          app.Redis,
	}
}

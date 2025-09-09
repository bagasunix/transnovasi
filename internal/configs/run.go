package configs

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/gofiber/fiber/v2"

	"github.com/bagasunix/transnovasi/internal/delivery/http"
	"github.com/bagasunix/transnovasi/pkg/env"
)

func Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := InitLogger()
	Cfg := env.InitConfig(ctx, logger)
	db := InitDB(ctx, Cfg, logger)
	redis := InitRedis(ctx, logger, Cfg)

	app := InitFiber(ctx, Cfg)
	setup := SetupApp(&setupApp{
		DB:    db,
		App:   app,
		Log:   logger,
		Cfg:   Cfg,
		Redis: redis,
	})
	httpHandler := http.InitHttpHandler(setup)

	errs := make(chan error)
	defer close(errs)
	go initCancel(errs)
	go initHttp(httpHandler, Cfg, errs)
	logger.Error().Msgf("exit: %v", <-errs)
}

func initCancel(errs chan error) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	errs <- fmt.Errorf("%s", <-c)
}

func initHttp(c *fiber.App, config *env.Cfg, errs chan error) {
	errs <- c.Listen(":" + strconv.Itoa(config.Server.Port))
}

package handlers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/bagasunix/transnovasi/internal/controllers"
)

func MakeCustHandler(controller *controllers.CustomerController, router fiber.Router, authMiddleware fiber.Handler) {
	router.Post("", authMiddleware, controller.Create)
	router.Get("", authMiddleware, controller.GetAllUser)
}

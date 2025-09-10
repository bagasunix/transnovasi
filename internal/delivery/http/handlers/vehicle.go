package handlers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/bagasunix/transnovasi/internal/controllers"
)

func MakeVehicleHandler(controller *controllers.VehicleController, router fiber.Router, authMiddleware fiber.Handler) {
	router.Get("", authMiddleware, controller.GetAllVehicle)
}

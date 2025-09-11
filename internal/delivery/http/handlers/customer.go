package handlers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/bagasunix/transnovasi/internal/controllers"
)

func MakeCustHandler(controller *controllers.CustomerController, router fiber.Router, authMiddleware fiber.Handler) {
	router.Post("", authMiddleware, controller.Create)
	router.Get("", authMiddleware, controller.GetAllCustomer)
	router.Get("/:id", authMiddleware, controller.ViewCustomer)
	router.Put("/:id", authMiddleware, controller.UpdateCustomer)
	router.Delete("/:id", authMiddleware, controller.DeleteCustomer)
	router.Get("/:id/vehicle", authMiddleware, controller.ViewCustomerWithVehicle)
}

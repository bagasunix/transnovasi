package handlers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/bagasunix/transnovasi/internal/controllers"
)

func MakeAuthHandler(controller *controllers.AuthController, router fiber.Router) {
	router.Post("customer", controller.LoginCustomer)
	router.Post("user", controller.LoginUser)
}

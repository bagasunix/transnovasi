package handlers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/bagasunix/transnovasi/internal/controllers"
)

func MakeAuthHandler(controller *controllers.AuthController, router fiber.Router, authMiddleware fiber.Handler) {
	router.Post("", controller.LoginUser)
	router.Post("register", controller.Register)
	router.Get("logout", authMiddleware, controller.Logout)
}

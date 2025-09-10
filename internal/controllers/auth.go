package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/phuslu/log"

	"github.com/bagasunix/transnovasi/internal/dtos/requests"
	"github.com/bagasunix/transnovasi/internal/dtos/responses"
	"github.com/bagasunix/transnovasi/internal/repositories"
	"github.com/bagasunix/transnovasi/internal/usecases"
)

type AuthController struct {
	usecase usecases.AuthUsecase
	logger  *log.Logger
	repo    repositories.Repositories
}

func NewAuthController(logger *log.Logger, repo repositories.Repositories, usecase usecases.AuthUsecase) *AuthController {
	return &AuthController{
		usecase: usecase,
		logger:  logger,
		repo:    repo,
	}
}

func (ac *AuthController) LoginUser(ctx *fiber.Ctx) error {
	req := new(requests.Login)
	var result responses.BaseResponse[*responses.ResponseLogin]
	if err := ctx.BodyParser(req); err != nil {
		result.Code = fiber.StatusBadRequest
		result.Message = "Invalid request"
		result.Errors = err.Error()
		return ctx.Status(fiber.StatusBadRequest).JSON(result)
	}

	result = ac.usecase.LoginUser(ctx.Context(), req)
	if result.Errors != "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(result)
	}

	return ctx.Status(result.Code).JSON(result)
}

func (ac *AuthController) Register(ctx *fiber.Ctx) error {
	req := new(requests.User)
	var result responses.BaseResponse[*responses.UserResponse]
	if err := ctx.BodyParser(req); err != nil {
		result.Code = fiber.StatusBadRequest
		result.Message = "Invalid request"
		result.Errors = err.Error()
		return ctx.Status(fiber.StatusBadRequest).JSON(result)
	}

	result = ac.usecase.Register(ctx.Context(), req)
	if result.Errors != "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(result)
	}

	return ctx.Status(result.Code).JSON(result)
}

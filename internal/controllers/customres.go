package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/phuslu/log"

	"github.com/bagasunix/transnovasi/internal/dtos/requests"
	"github.com/bagasunix/transnovasi/internal/dtos/responses"
	"github.com/bagasunix/transnovasi/internal/repositories"
	"github.com/bagasunix/transnovasi/internal/usecases"
)

type CustomerController struct {
	usecase usecases.CustomerUsecase
	logger  *log.Logger
	repo    repositories.Repositories
}

func NewCustomerController(logger *log.Logger, repo repositories.Repositories, usecase usecases.CustomerUsecase) *CustomerController {
	return &CustomerController{
		usecase: usecase,
		logger:  logger,
		repo:    repo,
	}
}

func (ac *CustomerController) Create(ctx *fiber.Ctx) error {
	req := new(requests.Customer)
	var result responses.BaseResponse[responses.CustomerResponse]
	if err := ctx.BodyParser(req); err != nil {
		result.Code = fiber.StatusBadRequest
		result.Message = "Invalid request"
		result.Errors = err.Error()
		return ctx.Status(fiber.StatusBadRequest).JSON(result)
	}

	result = ac.usecase.Create(ctx.Context(), req)
	if result.Errors != "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(result)
	}

	return ctx.Status(result.Code).JSON(result)
}

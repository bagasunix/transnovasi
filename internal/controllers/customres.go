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
	request := new(requests.CreateCustomer)
	var response responses.BaseResponse[responses.CustomerResponse]
	if err := ctx.BodyParser(request); err != nil {
		response.Code = fiber.StatusBadRequest
		response.Message = "Invalid request"
		response.Errors = err.Error()
		return ctx.Status(fiber.StatusBadRequest).JSON(response)
	}

	response = ac.usecase.Create(ctx.Context(), request)
	if response.Errors != "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(response)
	}

	return ctx.Status(response.Code).JSON(response)
}
func (c *CustomerController) GetAllCustomer(ctx *fiber.Ctx) error {
	request := new(requests.BaseRequest)
	var response responses.BaseResponse[[]responses.CustomerResponse]

	if err := ctx.QueryParser(request); err != nil {
		response.Code = fiber.StatusBadRequest
		response.Message = "Invalid request"
		response.Errors = err.Error()
		return ctx.Status(fiber.StatusBadRequest).JSON(response)
	}
	response = c.usecase.ListCustomer(ctx.Context(), request)
	if response.Errors != "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(response)
	}
	return ctx.Status(response.Code).JSON(response)
}
func (c *CustomerController) ViewCustomer(ctx *fiber.Ctx) error {
	request := new(requests.EntityId)
	var response responses.BaseResponse[*responses.CustomerResponse]

	id := ctx.Params("id")
	if id == "" {
		response.Code = fiber.StatusBadRequest
		response.Message = "Invalid request"
		response.Errors = "ID tidak ditemukan"
		return ctx.Status(fiber.StatusBadRequest).JSON(response)
	}

	request.Id = id
	response = c.usecase.ViewCustomer(ctx.Context(), request)
	if response.Errors != "" {
		return ctx.Status(response.Code).JSON(response)
	}

	return ctx.Status(response.Code).JSON(response)
}
func (c *CustomerController) UpdateCustomer(ctx *fiber.Ctx) error {
	request := new(requests.UpdateCustomer)
	var response responses.BaseResponse[*responses.CustomerResponse]

	id := ctx.Params("id")
	if id == "" {
		response.Code = fiber.StatusBadRequest
		response.Message = "Invalid request"
		response.Errors = "ID tidak ditemukan"
		return ctx.Status(fiber.StatusBadRequest).JSON(response)
	}

	if err := ctx.BodyParser(request); err != nil {
		response.Code = fiber.StatusBadRequest
		response.Message = "Invalid request"
		response.Errors = err.Error()
		return ctx.Status(fiber.StatusBadRequest).JSON(response)
	}
	request.ID = id
	response = c.usecase.UpdateCustomer(ctx.Context(), request)
	if response.Errors != "" {
		return ctx.Status(response.Code).JSON(response)
	}

	return ctx.Status(response.Code).JSON(response)
}

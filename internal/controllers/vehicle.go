package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/phuslu/log"

	"github.com/bagasunix/transnovasi/internal/dtos/requests"
	"github.com/bagasunix/transnovasi/internal/dtos/responses"
	"github.com/bagasunix/transnovasi/internal/repositories"
	"github.com/bagasunix/transnovasi/internal/usecases"
)

type VehicleController struct {
	usecase usecases.VehicleUsecase
	logger  *log.Logger
	repo    repositories.Repositories
}

func NewVehicleController(logger *log.Logger, repo repositories.Repositories, usecase usecases.VehicleUsecase) *VehicleController {
	return &VehicleController{
		usecase: usecase,
		logger:  logger,
		repo:    repo,
	}
}
func (c *VehicleController) GetAllVehicle(ctx *fiber.Ctx) error {
	request := new(requests.BaseVehicle)
	var response responses.BaseResponse[[]responses.VehicleResponse]

	if err := ctx.QueryParser(request); err != nil {
		response.Code = fiber.StatusBadRequest
		response.Message = "Invalid request"
		response.Errors = err.Error()
		return ctx.Status(fiber.StatusBadRequest).JSON(response)
	}
	response = c.usecase.ListVehicle(ctx.Context(), request)
	if response.Errors != "" {
		return ctx.Status(response.Code).JSON(response)
	}
	return ctx.Status(response.Code).JSON(response)
}

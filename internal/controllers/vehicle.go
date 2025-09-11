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

// GetAllVehicle godoc
// @Summary      Get all vehicles
// @Description  Retrieve a list of vehicles with optional filters
// @Tags         Vehicle
// @Accept       json
// @Produce      json
// @Param        search query string false "Search keyword for vehicle"
// @Param        limit query int false "Number of items per page" default(10)
// @Param        page query int false "Page number" default(1)
// @Success      200 {object} responses.BaseResponseVehicleList
// @Failure      400 {object} responses.BaseResponseSwagger
// @Router       /vehicles [get]
// @Security BearerAuth
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

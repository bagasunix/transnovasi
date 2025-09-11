package controllers

import (
	"strconv"

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

// Create godoc
// @Summary Membuat customer baru
// @Description Membuat data customer baru beserta kendaraan opsional
// @Tags Customer
// @Accept  json
// @Produce  json
// @Param   request body requests.CreateCustomer true "Customer Request"
// @Success 200 {object} responses.BaseResponseCustomer
// @Failure 400 {object} responses.BaseResponseSwagger
// @Failure 401 {object} responses.BaseResponseSwagger
// @Failure 409 {object} responses.BaseResponseSwagger
// @Router /customers [post]
// @Security BearerAuth
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
		return ctx.Status(response.Code).JSON(response)
	}

	return ctx.Status(response.Code).JSON(response)
}

// GetAllCustomer godoc
// @Summary Mendapatkan daftar customer
// @Description Menampilkan daftar customer dengan pagination dan optional search
// @Tags Customer
// @Accept  json
// @Produce  json
// @Param   page   query string false "Page number"
// @Param   limit  query string false "Limit per page"
// @Param   search query string false "Search by name/email"
// @Success 200 {object} responses.BaseResponseListCustomer
// @Failure 400 {object} responses.BaseResponseSwagger
// @Failure 401 {object} responses.BaseResponseSwagger
// @Router /customers [get]
// @Security BearerAuth
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
		return ctx.Status(response.Code).JSON(response)
	}
	return ctx.Status(response.Code).JSON(response)
}

// ViewCustomer godoc
// @Summary Mendapatkan detail customer
// @Description Menampilkan detail customer berdasarkan ID
// @Tags Customer
// @Accept  json
// @Produce  json
// @Param   id   path int true "Customer ID"
// @Success 200 {object} responses.BaseResponseCustomer
// @Failure 400 {object} responses.BaseResponseSwagger
// @Failure 404 {object} responses.BaseResponseSwagger
// @Router /customers/{id} [get]
// @Security BearerAuth
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
	if _, err := strconv.Atoi(id); err != nil {
		response.Code = fiber.StatusBadRequest
		response.Message = "Invalid request"
		response.Errors = "ID harus berupa angka"
		return ctx.Status(response.Code).JSON(response)
	}

	request.Id = id
	response = c.usecase.ViewCustomer(ctx.Context(), request)
	if response.Errors != "" {
		return ctx.Status(response.Code).JSON(response)
	}

	return ctx.Status(response.Code).JSON(response)
}

// UpdateCustomer godoc
// @Summary Melakukan update customer
// @Description Melakukan update customer berdasarkan ID
// @Tags Customer
// @Accept  json
// @Produce  json
// @Param   id   path int true "Customer ID"
// @Param   body body requests.UpdateCustomer true "Update Customer Request"
// @Success 200 {object} responses.BaseResponseCustomer
// @Failure 400 {object} responses.BaseResponseSwagger
// @Failure 404 {object} responses.BaseResponseSwagger
// @Router /customers/{id} [put]
// @Security BearerAuth
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

	if _, err := strconv.Atoi(id); err != nil {
		response.Code = fiber.StatusBadRequest
		response.Message = "Invalid request"
		response.Errors = "ID harus berupa angka"
		return ctx.Status(response.Code).JSON(response)
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

// DeleteCustomer godoc
// @Summary Melakukan delete customer
// @Description Melakukan delete customer berdasarkan ID
// @Tags Customer
// @Accept  json
// @Produce  json
// @Param   id   path int true "Customer ID"
// @Success 200 {object} responses.BaseResponseSwagger
// @Failure 400 {object} responses.BaseResponseSwagger
// @Failure 404 {object} responses.BaseResponseSwagger
// @Router /customers/{id} [delete]
// @Security BearerAuth
func (c *CustomerController) DeleteCustomer(ctx *fiber.Ctx) error {
	request := new(requests.EntityId)
	var response responses.BaseResponse[any]

	id := ctx.Params("id")
	if id == "" {
		response.Code = fiber.StatusBadRequest
		response.Message = "Invalid request"
		response.Errors = "ID tidak ditemukan"
		return ctx.Status(fiber.StatusBadRequest).JSON(response)
	}
	if _, err := strconv.Atoi(id); err != nil {
		response.Code = fiber.StatusBadRequest
		response.Message = "Invalid request"
		response.Errors = "ID harus berupa angka"
		return ctx.Status(response.Code).JSON(response)
	}

	request.Id = id
	response = c.usecase.DeleteCustomer(ctx.Context(), request)
	if response.Errors != "" {
		return ctx.Status(response.Code).JSON(response)
	}

	return ctx.Status(response.Code).JSON(response)
}

// ViewCustomerWithVehicle godoc
// @Summary Melihat kendaraan milik customer
// @Description Mengambil daftar kendaraan berdasarkan customer ID
// @Tags Customer
// @Accept  json
// @Produce  json
// @Param   id   path int true "Customer ID"
// @Success 200 {object} responses.BaseResponseVehicleList
// @Failure 400 {object} responses.BaseResponseSwagger
// @Failure 404 {object} responses.BaseResponseSwagger
// @Router /customers/{id}/vehicles [get]
// @Security BearerAuth
func (c *CustomerController) ViewCustomerWithVehicle(ctx *fiber.Ctx) error {
	request := new(requests.BaseVehicle)
	var response responses.BaseResponse[[]responses.VehicleResponse]

	id := ctx.Params("id")
	if id == "" {
		response.Code = fiber.StatusBadRequest
		response.Message = "Invalid request"
		response.Errors = "ID tidak ditemukan"
		return ctx.Status(fiber.StatusBadRequest).JSON(response)
	}
	if _, err := strconv.Atoi(id); err != nil {
		response.Code = fiber.StatusBadRequest
		response.Message = "Invalid request"
		response.Errors = "ID harus berupa angka"
		return ctx.Status(response.Code).JSON(response)
	}

	request.CustomerID = id
	response = c.usecase.ViewCustomerWithVehicle(ctx.Context(), request)
	if response.Errors != "" {
		return ctx.Status(response.Code).JSON(response)
	}

	return ctx.Status(response.Code).JSON(response)
}

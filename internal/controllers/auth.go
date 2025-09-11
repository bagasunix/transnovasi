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

// LoginUser godoc
// @Summary Login pengguna
// @Description Autentikasi pengguna menggunakan email dan password
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param   request body requests.Login true "Login Request"
// @Success 200 {object} responses.BaseResponseLogin
// @Failure 400 {object} responses.BaseResponseError
// @Router /auth/login [post]
func (ac *AuthController) LoginUser(ctx *fiber.Ctx) error {
	request := new(requests.Login)
	var result responses.BaseResponse[*responses.ResponseLogin]
	if err := ctx.BodyParser(request); err != nil {
		result.Code = fiber.StatusBadRequest
		result.Message = "Invalid request"
		result.Errors = err.Error()
		return ctx.Status(fiber.StatusBadRequest).JSON(result)
	}

	result = ac.usecase.LoginUser(ctx.Context(), request)
	if result.Errors != "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(result)
	}

	return ctx.Status(result.Code).JSON(result)
}

// Register godoc
// @Summary Daftar pengguna baru
// @Description Membuat akun pengguna baru
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param   request body requests.User true "Register Request"
// @Success 200 {object} responses.BaseResponseUser
// @Failure 400 {object} responses.BaseResponseError
// @Router /auth/register [post]
func (ac *AuthController) Register(ctx *fiber.Ctx) error {
	request := new(requests.User)
	var result responses.BaseResponse[*responses.UserResponse]
	if err := ctx.BodyParser(request); err != nil {
		result.Code = fiber.StatusBadRequest
		result.Message = "Invalid request"
		result.Errors = err.Error()
		return ctx.Status(fiber.StatusBadRequest).JSON(result)
	}

	result = ac.usecase.Register(ctx.Context(), request)
	if result.Errors != "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(result)
	}

	return ctx.Status(result.Code).JSON(result)
}

// Logout godoc
// @Summary Logout pengguna
// @Description Menghapus token pengguna dari sistem (invalidate session)
// @Tags Auth
// @Accept  json
// @Produce  json
// @Success 200 {object} responses.BaseResponseSwagger
// @Failure 400 {object} responses.BaseResponseError
// @Router /auth/logout [post]
// @Security BearerAuth
func (ac *AuthController) Logout(ctx *fiber.Ctx) error {
	result := ac.usecase.LogoutUser(ctx.Context())
	if result.Errors != "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(result)
	}
	return ctx.Status(result.Code).JSON(result)
}

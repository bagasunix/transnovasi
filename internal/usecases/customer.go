package usecases

import (
	"github.com/gofiber/fiber/v2"
	"github.com/phuslu/log"
	"gorm.io/gorm"

	"github.com/bagasunix/transnovasi/internal/dtos/requests"
	"github.com/bagasunix/transnovasi/internal/dtos/responses"
	"github.com/bagasunix/transnovasi/internal/repositories"
	"github.com/bagasunix/transnovasi/pkg/env"
)

type custUsecase struct {
	db     *gorm.DB
	cfg    *env.Cfg
	repo   repositories.Repositories
	logger *log.Logger
}

type CustomerUsecase interface {
	Create(ctx *fiber.Ctx, req *requests.Customer) (response responses.BaseResponse[responses.CustomerResponse])
}

func NewCustUsecase(logger *log.Logger, db *gorm.DB, cfg *env.Cfg, repo repositories.Repositories) CustomerUsecase {
	n := new(custUsecase)
	n.cfg = cfg
	n.db = db
	n.logger = logger
	n.repo = repo
	return n
}

func (c *custUsecase) Create(ctx *fiber.Ctx, req *requests.Customer) (response responses.BaseResponse[responses.CustomerResponse]) {
	return response
}

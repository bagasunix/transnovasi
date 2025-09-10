package usecases

import (
	"context"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/phuslu/log"
	"gorm.io/gorm"

	"github.com/bagasunix/transnovasi/internal/domains"
	"github.com/bagasunix/transnovasi/internal/dtos/requests"
	"github.com/bagasunix/transnovasi/internal/dtos/responses"
	"github.com/bagasunix/transnovasi/internal/middlewares"
	"github.com/bagasunix/transnovasi/internal/repositories"
	"github.com/bagasunix/transnovasi/pkg/env"
	"github.com/bagasunix/transnovasi/pkg/errors"
	"github.com/bagasunix/transnovasi/pkg/helpers"
)

type custUsecase struct {
	db     *gorm.DB
	cfg    *env.Cfg
	repo   repositories.Repositories
	logger *log.Logger
}

type CustomerUsecase interface {
	Create(ctx context.Context, request *requests.Customer) (response responses.BaseResponse[responses.CustomerResponse])
}

func NewCustUsecase(logger *log.Logger, db *gorm.DB, cfg *env.Cfg, repo repositories.Repositories) CustomerUsecase {
	n := new(custUsecase)
	n.cfg = cfg
	n.db = db
	n.logger = logger
	n.repo = repo
	return n
}

func (c *custUsecase) Create(ctx context.Context, request *requests.Customer) (response responses.BaseResponse[responses.CustomerResponse]) {
	authUser := middlewares.GetAuthClaimsFromContext(ctx)
	if request.Validate() != nil {
		response.Code = fiber.StatusBadRequest
		response.Message = "email atau password salah, silahkan coba lagi"
		response.Errors = request.Validate().Error()
		return response
	}

	checkEmail := c.repo.GetCustomer().GetOneByParams(ctx, map[string]any{"email": request.Email})
	if checkEmail.Value.Email == request.Email {
		response.Code = fiber.StatusConflict
		response.Message = "Email sudah terdaftar"
		response.Errors = "email " + errors.ERR_ALREADY_EXISTS
		return response
	}
	if checkEmail.Error != nil && !strings.Contains(checkEmail.Error.Error(), "not found") {
		response.Code = fiber.StatusConflict
		response.Message = "Validasi email invalid"
		response.Errors = checkEmail.Error.Error()
		return response
	}

	phone, err := helpers.ValidatePhone(request.Phone)
	if err != nil {
		response.Code = fiber.StatusConflict
		response.Message = "Validasi phone invalid"
		response.Errors = err.Error()
		return response
	}

	strUserID, _ := strconv.Atoi(authUser.ID)

	// Build customer
	customerBuild := &domains.Customer{
		Name:      request.Name,
		Email:     request.Email,
		Phone:     *phone,
		Address:   request.Address,
		IsActive:  1,
		CreatedBy: strUserID,
	}

	tx := c.repo.GetUser().GetConnection().(*gorm.DB).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err = c.repo.GetCustomer().CreateTx(ctx, tx, customerBuild); err != nil {
		response.Code = fiber.StatusConflict
		response.Message = "Gagal membuat pengguna"
		response.Errors = err.Error()
		return response
	}

	vehicles := make([]domains.Vehicle, 0, len(request.Vehicle))
	if len(request.Vehicle) != 0 {
		platSet := make(map[string]struct{}, len(request.Vehicle))
		for _, v := range request.Vehicle {
			// Cek duplikat dalam request
			if _, exists := platSet[v.PlateNo]; exists {
				response.Code = fiber.StatusBadRequest
				response.Message = "Nomor plat duplikat dalam request"
				response.Errors = "Plat " + v.PlateNo + " sudah ada di request"
				return response
			}
			platSet[v.PlateNo] = struct{}{}

			checkPlat := c.repo.GetVehicle().GetOneByParams(ctx, map[string]interface{}{"plate_no": v.PlateNo})
			if checkPlat.Value.PlateNo == v.PlateNo {
				response.Code = fiber.StatusConflict
				response.Message = "Kendaraan sudah terdaftar"
				response.Errors = "vehicle " + errors.ERR_ALREADY_EXISTS
				return response
			}
			if checkPlat.Error != nil && !strings.Contains(checkPlat.Error.Error(), "not found") {
				response.Code = fiber.StatusConflict
				response.Message = "Validasi vehicle invalid"
				response.Errors = checkPlat.Error.Error()
				return response
			}

			intYear, convErr := strconv.Atoi(v.Year)
			if convErr != nil {
				return responses.BaseResponse[responses.CustomerResponse]{
					Code:    fiber.StatusBadRequest,
					Message: "Tahun kendaraan tidak valid",
					Errors:  convErr.Error(),
				}
			}

			vehicles = append(vehicles, domains.Vehicle{
				Brand:      v.Brand,
				Color:      v.Color,
				CustomerID: customerBuild.ID,
				FuelType:   v.FuelType,
				MaxSpeed:   v.MaxSpeed,
				Model:      v.Model,
				PlateNo:    v.PlateNo,
				Year:       intYear,
				CreatedBy:  strUserID,
			})
		}
	}

	if len(vehicles) != 0 {
		if err := c.repo.GetVehicle().CreateBulkTx(ctx, tx, vehicles); err != nil {
			tx.Rollback()
			return responses.BaseResponse[responses.CustomerResponse]{
				Code:    fiber.StatusConflict,
				Message: "Gagal membuat pengguna",
				Errors:  err.Error(),
			}
		}
	}

	if err = tx.Commit().Error; err != nil {
		response.Code = fiber.StatusConflict
		response.Message = "Gagal membuat pengguna"
		response.Errors = err.Error()
		return response
	}

	// Build response
	resBuild := &responses.CustomerResponse{
		ID:       int64(customerBuild.ID),
		Name:     customerBuild.Name,
		Phone:    customerBuild.Phone,
		Email:    customerBuild.Email,
		Address:  customerBuild.Address,
		IsActive: strconv.Itoa(customerBuild.IsActive),
	}

	if len(vehicles) > 0 {
		resBuild.Vehicle = make([]responses.VehicleResponse, 0, len(vehicles))
		for _, v := range vehicles {
			resBuild.Vehicle = append(resBuild.Vehicle, responses.VehicleResponse{
				ID:       strconv.Itoa(v.ID),
				Brand:    v.Brand,
				Color:    v.Color,
				FuelType: v.FuelType,
				MaxSpeed: v.MaxSpeed,
				Model:    v.Model,
				PlateNo:  v.PlateNo,
				Year:     v.Year,
				IsActive: v.IsActive,
			})
		}
	}

	response.Code = fiber.StatusOK
	response.Message = "Sukses mendaftar"
	response.Data = resBuild
	return response
}

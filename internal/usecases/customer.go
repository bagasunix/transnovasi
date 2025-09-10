package usecases

import (
	errs "errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/phuslu/log"
	"gorm.io/gorm"

	"github.com/bagasunix/transnovasi/internal/domains"
	"github.com/bagasunix/transnovasi/internal/dtos/requests"
	"github.com/bagasunix/transnovasi/internal/dtos/responses"
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
	Create(ctx *fiber.Ctx, request *requests.Customer) (response responses.BaseResponse[responses.CustomerResponse])
}

func NewCustUsecase(logger *log.Logger, db *gorm.DB, cfg *env.Cfg, repo repositories.Repositories) CustomerUsecase {
	n := new(custUsecase)
	n.cfg = cfg
	n.db = db
	n.logger = logger
	n.repo = repo
	return n
}

func (c *custUsecase) Create(ctx *fiber.Ctx, request *requests.Customer) (response responses.BaseResponse[responses.CustomerResponse]) {
	if request.Validate() != nil {
		response.Code = fiber.StatusBadRequest
		response.Message = "email atau password salah, silahkan coba lagi"
		response.Errors = request.Validate().Error()
		return response
	}

	checkEmail := c.repo.GetCustomer().GetOneByParams(ctx.Context(), map[string]interface{}{"email": request.Email})
	if checkEmail.Value.Email == request.Email {
		response.Code = fiber.StatusConflict
		response.Message = "Email sudah terdaftar"
		response.Errors = "email " + errors.ERR_ALREADY_EXISTS
		return response
	}
	if checkEmail.Error != nil && !errs.Is(checkEmail.Error, gorm.ErrRecordNotFound) {
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

	customerBuild := new(domains.Customer)
	customerBuild.Name = request.Name
	customerBuild.Email = request.Email
	customerBuild.Phone = *phone
	customerBuild.Address = request.Address
	customerBuild.IsActive = 1

	var vehicles []domains.Vehicle
	if len(request.Vehicle) != 0 {
		for _, v := range request.Vehicle {
			checkPlat := c.repo.GetVehicle().GetOneByParams(ctx.Context(), map[string]interface{}{"plat_no": v.PlateNo})
			if checkPlat.Value.PlateNo == v.PlateNo {
				response.Code = fiber.StatusConflict
				response.Message = "Kendaraan sudah terdaftar"
				response.Errors = "vehicle " + errors.ERR_ALREADY_EXISTS
				return response
			}
			if checkPlat.Error != nil && !errs.Is(checkPlat.Error, gorm.ErrRecordNotFound) {
				response.Code = fiber.StatusConflict
				response.Message = "Validasi vehicle invalid"
				response.Errors = checkPlat.Error.Error()
				return response
			}

			intYear, _ := strconv.Atoi(v.Year)

			vehicleBuild := new(domains.Vehicle)
			vehicleBuild.Brand = v.Brand
			vehicleBuild.Color = v.Color
			vehicleBuild.CustomerID = customerBuild.ID
			vehicleBuild.FuelType = v.FuelType
			vehicleBuild.MaxSpeed = v.MaxSpeed
			vehicleBuild.Model = v.Model
			vehicleBuild.PlateNo = v.PlateNo
			vehicleBuild.Year = intYear

			vehicles = append(vehicles, *vehicleBuild)
		}
	}

	tx := c.repo.GetUser().GetConnection().(*gorm.DB).Begin()

	if err = c.repo.GetCustomer().CreateTx(ctx.Context(), tx, customerBuild); err != nil {
		response.Code = fiber.StatusConflict
		response.Message = "Gagal membuat pengguna"
		response.Errors = err.Error()
		return response
	}

	if err = c.repo.GetVehicle().CreateBulkTx(ctx.Context(), tx, vehicles); err != nil {
		response.Code = fiber.StatusConflict
		response.Message = "Gagal membuat pengguna"
		response.Errors = err.Error()
		return response
	}

	if err = tx.Commit().Error; err != nil {
		response.Code = fiber.StatusConflict
		response.Message = "Gagal membuat pengguna"
		response.Errors = err.Error()
		return response
	}

	resBuild := new(responses.CustomerResponse)
	resBuild.ID = int64(customerBuild.ID)
	resBuild.Name = customerBuild.Name
	resBuild.Phone = customerBuild.Phone
	resBuild.Email = customerBuild.Email
	resBuild.Address = customerBuild.Address
	resBuild.IsActive = strconv.Itoa(customerBuild.IsActive)

	if len(vehicles) != 0 {
		for _, v := range vehicles {
			vehicle := new(responses.VehicleResponse)
			vehicle.Brand = v.Brand
			vehicle.Color = v.Color
			vehicle.FuelType = v.FuelType
			vehicle.MaxSpeed = v.MaxSpeed
			vehicle.Model = v.Model
			vehicle.PlateNo = v.PlateNo
			vehicle.Year = v.Year

			resBuild.Vehicle = append(resBuild.Vehicle, *vehicle)
		}
	}

	response.Code = fiber.StatusOK
	response.Message = "Sukses mendaftar"
	response.Data = resBuild
	return response
}

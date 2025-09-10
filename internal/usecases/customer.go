package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
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
	redis  *redis.Client
}

type CustomerUsecase interface {
	Create(ctx context.Context, request *requests.CreateCustomer) (response responses.BaseResponse[responses.CustomerResponse])
	ListCustomer(ctx context.Context, request *requests.BaseRequest) (response responses.BaseResponse[[]responses.CustomerResponse])
	ViewCustomer(ctx context.Context, request *requests.EntityId) (response responses.BaseResponse[*responses.CustomerResponse])
	UpdateCustomer(ctx context.Context, request *requests.UpdateCustomer) (response responses.BaseResponse[*responses.CustomerResponse])
}

func NewCustUsecase(logger *log.Logger, db *gorm.DB, cfg *env.Cfg, repo repositories.Repositories, redis *redis.Client) CustomerUsecase {
	n := new(custUsecase)
	n.cfg = cfg
	n.db = db
	n.logger = logger
	n.repo = repo
	n.redis = redis
	return n
}

func (c *custUsecase) Create(ctx context.Context, request *requests.CreateCustomer) (response responses.BaseResponse[responses.CustomerResponse]) {
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
		response.Message = "Gagal membuat pelanggan"
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
				Message: "Gagal membuat pelanggan",
				Errors:  err.Error(),
			}
		}
	}

	if err = tx.Commit().Error; err != nil {
		response.Code = fiber.StatusConflict
		response.Message = "Gagal membuat pelanggan"
		response.Errors = err.Error()
		return response
	}

	// Hapus cache Redis
	// hapus semua key yang dimulai dengan "customers:"
	keys, _ := c.redis.Keys(ctx, "customers:*").Result()
	if len(keys) > 0 {
		c.redis.Del(ctx, keys...)
	}

	// Build response
	resBuild := &responses.CustomerResponse{
		ID:       strconv.Itoa(customerBuild.ID),
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
func (c *custUsecase) ListCustomer(ctx context.Context, request *requests.BaseRequest) (response responses.BaseResponse[[]responses.CustomerResponse]) {
	if err := request.Validate(); err != nil {
		response.Code = 400
		response.Message = "Error Validasi"
		response.Errors = err.Error()
		return response
	}
	intPage, _ := strconv.Atoi(request.Page)
	intLimit, _ := strconv.Atoi(request.Limit)

	offset, limit := helpers.CalculateOffsetAndLimit(intPage, intLimit)

	// --- Redis key berdasarkan search, page, limit
	cacheKey := fmt.Sprintf("customers:search=%s:page=%d:limit=%d", request.Search, intPage, intLimit)

	var custResponse []responses.CustomerResponse

	// Cek Redis
	val, err := c.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		// Cache hit, unmarshal langsung
		if err := json.Unmarshal([]byte(val), &custResponse); err == nil {
			// Hit, kembalikan response dengan paging
			totalItems, _ := c.redis.Get(ctx, "customers:count:"+request.Search).Int() // optional cache count
			totalPages := (totalItems + limit - 1) / limit
			response.Data = &custResponse
			response.Paging = &responses.PageMetadata{
				Page:      intPage,
				Size:      limit,
				TotalItem: totalItems,
				TotalPage: totalPages,
			}
			response.Message = "Inquiry pelanggan berhasil"
			response.Code = fiber.StatusOK
			return response
		}
	}

	//  Ambil dari DB
	resCust := c.repo.GetCustomer().GetAll(ctx, limit, offset, request.Search)
	if resCust.Error != nil {
		response.Code = 400
		response.Message = "Gagal menarik data"
		response.Errors = resCust.Error.Error()
		return response
	}
	// Calculate total items and total pages
	totalItems, err := c.repo.GetCustomer().CountCustomer(ctx, request.Search)
	if err != nil {
		response.Code = 400
		response.Message = "Gagal menarik data"
		response.Errors = err.Error()
		return response
	}
	totalPages := (totalItems + limit - 1) / limit

	// Map ke response
	custResponse = make([]responses.CustomerResponse, 0, len(resCust.Value))
	for _, v := range resCust.Value {
		custResponse = append(custResponse, responses.CustomerResponse{
			ID:       strconv.Itoa(v.ID),
			Name:     v.Name,
			Email:    v.Email,
			Phone:    v.Phone,
			Address:  v.Address,
			IsActive: strconv.Itoa(v.IsActive),
		})
	}

	// Simpan ke Redis
	data, _ := json.Marshal(custResponse)
	c.redis.Set(ctx, cacheKey, data, 5*time.Minute)
	c.redis.Set(ctx, "customers:count:"+request.Search, totalItems, 5*time.Minute)

	response.Data = &custResponse
	response.Paging = &responses.PageMetadata{
		Page:      intPage,
		Size:      limit,
		TotalItem: totalItems,
		TotalPage: totalPages,
	}
	response.Message = "Inquiry pelanggan berhasil"
	response.Code = fiber.StatusOK
	return response
}
func (c *custUsecase) ViewCustomer(ctx context.Context, request *requests.EntityId) (response responses.BaseResponse[*responses.CustomerResponse]) {
	if err := request.Validate(); err != nil {
		response.Code = fiber.StatusBadRequest
		response.Message = "Validasi error"
		response.Errors = err.Error()
		return response
	}

	paramID, _ := strconv.Atoi(request.Id.(string))
	// --- Redis key berdasarkan customer ID
	cacheKey := fmt.Sprintf("customers:%d", paramID)

	//  Cek Redis dulu
	resCust := new(responses.CustomerResponse)
	val, err := c.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		// Cache hit, unmarshal JSON
		if err := json.Unmarshal([]byte(val), &resCust); err == nil {
			response.Data = &resCust
			response.Message = "Pelanggan ditemukan"
			response.Code = fiber.StatusOK
			return response
		}
	}

	checkCustomer := c.repo.GetCustomer().GetOneByParams(ctx, map[string]any{"id": paramID})
	if checkCustomer.Error != nil {
		response.Code = fiber.StatusNotFound
		response.Message = "Pelanggan tidak ditemukan"
		response.Errors = checkCustomer.Error.Error()
		return response
	}

	resCust.ID = strconv.Itoa(checkCustomer.Value.ID)
	resCust.Name = checkCustomer.Value.Name
	resCust.Email = checkCustomer.Value.Email
	resCust.Phone = checkCustomer.Value.Phone
	resCust.Address = checkCustomer.Value.Address
	resCust.IsActive = strconv.Itoa(checkCustomer.Value.IsActive)

	if len(checkCustomer.Value.Vehicles) != 0 {
		resCust.Vehicle = make([]responses.VehicleResponse, 0, len(checkCustomer.Value.Vehicles))
		for _, v := range checkCustomer.Value.Vehicles {
			resCust.Vehicle = append(resCust.Vehicle, responses.VehicleResponse{
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

	// Simpan ke Redis dengan expire 5 menit
	data, _ := json.Marshal(resCust)
	c.redis.Set(ctx, cacheKey, data, 5*time.Minute)

	response.Data = &resCust
	response.Message = "Pelanggan ditemukan"
	response.Code = fiber.StatusOK

	return response
}
func (c *custUsecase) UpdateCustomer(ctx context.Context, request *requests.UpdateCustomer) (response responses.BaseResponse[*responses.CustomerResponse]) {
	if err := request.Validate(); err != nil {
		response.Code = fiber.StatusBadRequest
		response.Message = "Validasi error"
		response.Errors = request.Validate().Error()
		return response
	}

	checkCust := c.repo.GetCustomer().GetOneByParams(ctx, map[string]any{"id": request.ID})
	if checkCust.Error != nil {
		response.Code = fiber.StatusConflict
		response.Message = "Pelanggan tidak ditemukan"
		response.Errors = checkCust.Error.Error()
		return response
	}

	mCustt := new(domains.Customer)
	mCustt.Name = request.Name
	mCustt.Phone = request.Phone
	mCustt.Address = request.Address

	// --- Redis key berdasarkan customer ID
	intCustID, _ := strconv.Atoi(request.ID)
	cacheKey := fmt.Sprintf("customers:%d", intCustID)

	if err := c.repo.GetCustomer().Updates(ctx, intCustID, mCustt); err != nil {
		response.Code = fiber.StatusBadRequest
		response.Message = "Gagal memperbarui pelanggan"
		response.Errors = err.Error()
		return response
	}

	resCust := new(responses.CustomerResponse)
	resCust.ID = strconv.Itoa(mCustt.ID)
	resCust.Name = mCustt.Name
	resCust.Email = mCustt.Email
	resCust.Phone = mCustt.Phone
	resCust.Address = mCustt.Address
	resCust.IsActive = strconv.Itoa(mCustt.IsActive)

	if len(checkCust.Value.Vehicles) != 0 {
		resCust.Vehicle = make([]responses.VehicleResponse, 0, len(checkCust.Value.Vehicles))
		for _, v := range checkCust.Value.Vehicles {
			resCust.Vehicle = append(resCust.Vehicle, responses.VehicleResponse{
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

	// Hapus cache detail
	c.redis.Del(ctx, cacheKey)

	// 4. Hapus cache list (opsional)
	listKeys, _ := c.redis.Keys(ctx, "customers:search=*").Result()
	if len(listKeys) > 0 {
		c.redis.Del(ctx, listKeys...)
	}

	response.Data = &resCust
	response.Code = fiber.StatusOK
	response.Message = "Berhsail memperbarui pelanggan"
	return response
}

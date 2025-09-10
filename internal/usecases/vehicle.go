package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/phuslu/log"
	"gorm.io/gorm"

	"github.com/bagasunix/transnovasi/internal/dtos/requests"
	"github.com/bagasunix/transnovasi/internal/dtos/responses"
	"github.com/bagasunix/transnovasi/internal/repositories"
	"github.com/bagasunix/transnovasi/pkg/env"
	"github.com/bagasunix/transnovasi/pkg/helpers"
)

type vehicleUsecase struct {
	db     *gorm.DB
	cfg    *env.Cfg
	repo   repositories.Repositories
	logger *log.Logger
	redis  *redis.Client
}

type VehicleUsecase interface {
	Create(ctx context.Context, request *[]requests.CreateVehicle) (response responses.BaseResponse[responses.VehicleResponse])
	ListVehicle(ctx context.Context, request *requests.BaseVehicle) (response responses.BaseResponse[[]responses.VehicleResponse])
	ViewVehicle(ctx context.Context, request *requests.EntityId) (response responses.BaseResponse[*responses.VehicleResponse])
	DeleteVehicle(ctx context.Context, request *requests.EntityId) (response responses.BaseResponse[any])
	// UpdateVehicle(ctx context.Context, request *requests.UpdateCustomer) (response responses.BaseResponse[*responses.VehicleResponse])
}

func NewVehicleUsecase(logger *log.Logger, db *gorm.DB, cfg *env.Cfg, repo repositories.Repositories, redis *redis.Client) VehicleUsecase {
	n := new(vehicleUsecase)
	n.cfg = cfg
	n.db = db
	n.logger = logger
	n.repo = repo
	n.redis = redis
	return n
}

func (v *vehicleUsecase) Create(ctx context.Context, request *[]requests.CreateVehicle) (response responses.BaseResponse[responses.VehicleResponse]) {
	return response
}
func (v *vehicleUsecase) ListVehicle(ctx context.Context, request *requests.BaseVehicle) (response responses.BaseResponse[[]responses.VehicleResponse]) {
	var cacheKey, countKey string
	if err := request.Validate(); err != nil {
		response.Code = 400
		response.Message = "Error Validasi"
		response.Errors = err.Error()
		return response
	}

	intPage, _ := strconv.Atoi(request.Page)
	intLimit, _ := strconv.Atoi(request.Limit)
	intCustomerID, _ := strconv.Atoi(request.CustomerID)

	offset, limit := helpers.CalculateOffsetAndLimit(intPage, intLimit)

	// --- Redis key berdasarkan search, page, limit
	if intCustomerID != 0 {
		// Ambil kendaraan by customer
		cacheKey = fmt.Sprintf("vehicles:customer:%s:search=%s:page=%d:limit=%d", request.CustomerID, request.Search, intPage, intLimit)
		countKey = fmt.Sprintf("vehicles:customer:%s:count:search=%s", request.CustomerID, request.Search)
	} else {
		// Ambil semua kendaraan (global)
		cacheKey = fmt.Sprintf("vehicles:search=%s:page=%d:limit=%d", request.Search, intPage, intLimit)
		countKey = fmt.Sprintf("vehicles:count:search=%s", request.Search)
	}

	var vehicleResponse []responses.VehicleResponse
	// Cek Redis
	val, err := v.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		// Cache hit, unmarshal langsung
		if err := json.Unmarshal([]byte(val), &vehicleResponse); err == nil {
			// Hit, kembalikan response dengan paging
			totalItems, _ := v.redis.Get(ctx, countKey).Int() // optional cache count
			totalPages := (totalItems + limit - 1) / limit
			response.Data = &vehicleResponse
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
	resCust := v.repo.GetVehicle().GetAllByCustomerID(ctx, intCustomerID, limit, offset, request.Search)
	if resCust.Error != nil {
		response.Code = 400
		response.Message = "Gagal menarik data"
		response.Errors = resCust.Error.Error()
		return response
	}
	// Calculate total items and total pages
	totalItems, err := v.repo.GetVehicle().CountVehicle(ctx, request.Search)
	if err != nil {
		response.Code = 400
		response.Message = "Gagal menarik data"
		response.Errors = err.Error()
		return response
	}
	totalPages := (totalItems + limit - 1) / limit

	// Map ke response
	vehicleResponse = make([]responses.VehicleResponse, 0, len(resCust.Value))
	for _, v := range resCust.Value {
		vehicleResponse = append(vehicleResponse, responses.VehicleResponse{
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

	// Simpan ke Redis
	data, _ := json.Marshal(vehicleResponse)
	v.redis.Set(ctx, cacheKey, data, 5*time.Minute)
	v.redis.Set(ctx, countKey, totalItems, 5*time.Minute)

	response.Data = &vehicleResponse
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
func (v *vehicleUsecase) ViewVehicle(ctx context.Context, request *requests.EntityId) (response responses.BaseResponse[*responses.VehicleResponse]) {
	return response
}
func (v *vehicleUsecase) DeleteVehicle(ctx context.Context, request *requests.EntityId) (response responses.BaseResponse[any]) {
	return response
}

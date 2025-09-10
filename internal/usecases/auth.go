package usecases

import (
	"context"
	errs "errors"
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
	"github.com/bagasunix/transnovasi/pkg/hash"
	"github.com/bagasunix/transnovasi/pkg/jwt"
)

type authUsecase struct {
	db     *gorm.DB
	redis  *redis.Client
	cfg    *env.Cfg
	repo   repositories.Repositories
	logger *log.Logger
}

type AuthUsecase interface {
	LoginUser(ctx context.Context, request *requests.Login) (response responses.BaseResponse[*responses.ResponseLogin])
	LoginCustomer(ctx context.Context, request *requests.Login) (response responses.BaseResponse[*responses.ResponseLogin])
}

func NewAuthUsecase(logger *log.Logger, db *gorm.DB, cfg *env.Cfg, repo repositories.Repositories, redis *redis.Client) AuthUsecase {
	n := new(authUsecase)
	n.cfg = cfg
	n.db = db
	n.logger = logger
	n.redis = redis
	n.repo = repo
	return n
}

func (au *authUsecase) LoginUser(ctx context.Context, request *requests.Login) (response responses.BaseResponse[*responses.ResponseLogin]) {
	if request.Validate() != nil {
		response.Code = fiber.StatusBadRequest
		response.Message = "email atau password salah, silahkan coba lagi"
		response.Errors = request.Validate().Error()
		return response
	}
	// Check Login User
	checkUser := au.repo.GetUser().GetOneByParams(ctx, map[string]interface{}{"email": request.Email})
	if len(checkUser.Value.Email) == 0 || errs.Is(checkUser.Error, gorm.ErrRecordNotFound) {
		response.Code = fiber.StatusNotFound
		response.Message = "email atau password salah, silahkan coba lagi"
		response.Errors = fmt.Sprintf("%s %s", request.Email, gorm.ErrRecordNotFound)
		return response
	}

	if checkUser.Error != nil && !errs.Is(checkUser.Error, gorm.ErrRecordNotFound) {
		response.Code = fiber.StatusNotFound
		response.Message = checkUser.Error.Error()
		response.Errors = checkUser.Error.Error()
		return response
	}

	if !hash.ComparePasswords(checkUser.Value.Password, []byte(request.Password)) {
		response.Code = fiber.StatusNotFound
		response.Message = "email atau password salah, silahkan coba lagi"
		response.Errors = "email atau password salah, silahkan coba lagi"
		return response
	}

	userBuild := responses.UserResponse{}
	userBuild.ID = int64(checkUser.Value.ID)
	userBuild.Name = checkUser.Value.Name
	userBuild.Sex = checkUser.Value.Sex
	userBuild.Email = checkUser.Value.Email
	userBuild.RoleID = checkUser.Value.RoleID
	userBuild.IsActive = int16(checkUser.Value.IsActive)

	clm := new(jwt.Claims)
	clm.User = &userBuild
	clm.ExpiresAt = time.Now().Add(24 * time.Hour).Unix()

	token, err := jwt.GenerateToken(au.cfg.Server.Token.JWTKey, *clm)
	if err != nil {
		response.Code = fiber.StatusConflict
		response.Message = "email atau password salah, silahkan coba lagi"
		response.Errors = err.Error()
		return response
	}

	redisKey := "auth_user:token:" + strconv.Itoa(int(userBuild.ID))
	err = au.redis.Set(ctx, redisKey, token, 24*time.Hour).Err()
	if err != nil {
		response.Code = fiber.StatusConflict
		response.Message = "email atau password salah, silahkan coba lagi"
		response.Errors = err.Error()
		return response
	}

	resBuild := new(responses.ResponseLogin)
	resBuild.ID = strconv.Itoa(int(userBuild.ID))
	resBuild.Token = token

	response.Data = &resBuild
	response.Code = fiber.StatusOK
	response.Message = "Pengguna berhasil masuk"
	return response
}

func (au *authUsecase) LoginCustomer(ctx context.Context, request *requests.Login) (response responses.BaseResponse[*responses.ResponseLogin]) {
	if request.Validate() != nil {
		response.Code = fiber.StatusBadRequest
		response.Message = "Validasi error"
		response.Message = "email atau password salah, silahkan coba lagi"
		response.Errors = request.Validate().Error()
		return response
	}

	// Check Login Customer
	checkCust := au.repo.GetCustomer().GetOneByParams(ctx, map[string]interface{}{"email": request.Email})
	if len(checkCust.Value.Email) == 0 || errs.Is(checkCust.Error, gorm.ErrRecordNotFound) {
		response.Code = fiber.StatusNotFound
		response.Message = "Email tidak ditemukan"
		response.Message = "email atau password salah, silahkan coba lagi"
		response.Errors = fmt.Sprintf("%s %s", request.Email, gorm.ErrRecordNotFound)
		return response
	}

	if checkCust.Error != nil && !errs.Is(checkCust.Error, gorm.ErrRecordNotFound) {
		response.Code = fiber.StatusNotFound
		response.Message = checkCust.Error.Error()
		response.Message = "email atau password salah, silahkan coba lagi"
		response.Errors = checkCust.Error.Error()
		return response
	}

	if !hash.ComparePasswords(checkCust.Value.Password, []byte(request.Password)) {
		response.Code = fiber.StatusNotFound
		response.Message = "username and password salah"
		response.Message = "email atau password salah, silahkan coba lagi"
		response.Errors = "email atau password salah, silahkan coba lagi"
		return response
	}

	custBuild := responses.CustomerResponse{}
	custBuild.ID = int64(checkCust.Value.ID)
	custBuild.Name = checkCust.Value.Name
	custBuild.Sex = checkCust.Value.Sex
	custBuild.Email = checkCust.Value.Email
	custBuild.Address = checkCust.Value.Address
	custBuild.RoleID = checkCust.Value.RoleID
	custBuild.IsActive = int16(checkCust.Value.IsActive)

	clm := new(jwt.Claims)
	clm.Customer = &custBuild
	clm.ExpiresAt = time.Now().Add(24 * time.Hour).Unix()

	token, err := jwt.GenerateToken(au.cfg.Server.Token.JWTKey, *clm)
	if err != nil {
		response.Code = fiber.StatusConflict
		response.Message = "Gagal membuat token"
		response.Message = "email atau password salah, silahkan coba lagi"
		response.Errors = err.Error()
		return response
	}

	redisKey := "auth_customer:token:" + strconv.Itoa(int(custBuild.ID))
	err = au.redis.Set(ctx, redisKey, token, 24*time.Hour).Err()
	if err != nil {
		response.Code = fiber.StatusConflict
		response.Message = "Gagal menyimpan token di Redis"
		response.Message = "email atau password salah, silahkan coba lagi"
		response.Errors = err.Error()
		return response
	}

	resBuild := new(responses.ResponseLogin)
	resBuild.ID = strconv.Itoa(int(custBuild.ID))
	resBuild.Token = token

	response.Data = &resBuild
	response.Code = fiber.StatusOK
	response.Message = "Pengguna berhasil masuk"
	return response
}

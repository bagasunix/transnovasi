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
	"github.com/bagasunix/transnovasi/pkg/errors"
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
	LoginUser(ctx context.Context, request *requests.Login) (resonse responses.BaseResponse[*responses.ResponseLogin])
	LoginCustomer(ctx context.Context, request *requests.Login) (resonse responses.BaseResponse[*responses.ResponseLogin])
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

func (au *authUsecase) LoginUser(ctx context.Context, request *requests.Login) (resonse responses.BaseResponse[*responses.ResponseLogin]) {
	if request.Validate() != nil {
		resonse.Code = fiber.StatusBadRequest
		resonse.Message = "email atau password salah, silahkan coba lagi"
		resonse.Errors = request.Validate()
		return resonse
	}
	// Check Login User
	checkUser := au.repo.GetUser().GetOneByParams(ctx, map[string]interface{}{"email": request.Email})
	if len(checkUser.Value.Email) == 0 || errs.Is(checkUser.Error, gorm.ErrRecordNotFound) {
		resonse.Code = fiber.StatusNotFound
		resonse.Message = "email atau password salah, silahkan coba lagi"
		resonse.Errors = errors.CustomError(fmt.Sprintf("%s %s", request.Email, gorm.ErrRecordNotFound))
		return resonse
	}

	if checkUser.Error != nil && !errs.Is(checkUser.Error, gorm.ErrRecordNotFound) {
		resonse.Code = fiber.StatusNotFound
		resonse.Message = checkUser.Error.Error()
		resonse.Errors = checkUser.Error
		return resonse
	}

	if !hash.ComparePasswords(checkUser.Value.Password, []byte(request.Password)) {
		resonse.Code = fiber.StatusNotFound
		resonse.Message = "email atau password salah, silahkan coba lagi"
		resonse.Errors = errors.ErrInvalidAttributes("email atau password salah, silahkan coba lagi")
		return resonse
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
	clm.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()

	token, err := jwt.GenerateToken(au.cfg.Server.Token.JWTKey, *clm)
	if err != nil {
		resonse.Code = fiber.StatusConflict
		resonse.Message = "email atau password salah, silahkan coba lagi"
		resonse.Errors = err
		return resonse
	}

	redisKey := "authUser:token:" + strconv.Itoa(int(userBuild.ID))
	err = au.redis.Set(ctx, redisKey, token, time.Hour).Err()
	if err != nil {
		resonse.Code = fiber.StatusConflict
		resonse.Message = "email atau password salah, silahkan coba lagi"
		resonse.Errors = err
		return resonse
	}

	resBuild := new(responses.ResponseLogin)
	resBuild.ID = strconv.Itoa(int(userBuild.ID))
	resBuild.Token = token

	resonse.Data = &resBuild
	resonse.Code = fiber.StatusOK
	resonse.Message = "Pengguna berhasil masuk"
	return resonse
}

func (au *authUsecase) LoginCustomer(ctx context.Context, request *requests.Login) (resonse responses.BaseResponse[*responses.ResponseLogin]) {
	if request.Validate() != nil {
		resonse.Code = fiber.StatusBadRequest
		resonse.Message = "Validasi error"
		resonse.Message = "email atau password salah, silahkan coba lagi"
		resonse.Errors = request.Validate()
		return resonse
	}

	// Check Login Customer
	checkCust := au.repo.GetCustomer().GetOneByParams(ctx, map[string]interface{}{"email": request.Email})
	if len(checkCust.Value.Email) == 0 || errs.Is(checkCust.Error, gorm.ErrRecordNotFound) {
		resonse.Code = fiber.StatusNotFound
		resonse.Message = "Email tidak ditemukan"
		resonse.Message = "email atau password salah, silahkan coba lagi"
		resonse.Errors = errors.CustomError("email " + errors.ERR_NOT_FOUND)
		return resonse
	}

	if checkCust.Error != nil && !errs.Is(checkCust.Error, gorm.ErrRecordNotFound) {
		resonse.Code = fiber.StatusNotFound
		resonse.Message = checkCust.Error.Error()
		resonse.Message = "email atau password salah, silahkan coba lagi"
		resonse.Errors = checkCust.Error
		return resonse
	}

	if !hash.ComparePasswords(checkCust.Value.Password, []byte(request.Password)) {
		resonse.Code = fiber.StatusNotFound
		resonse.Message = "username and password salah"
		resonse.Message = "email atau password salah, silahkan coba lagi"
		resonse.Errors = errors.ErrInvalidAttributes("username and password")
		return resonse
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
	clm.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()

	token, err := jwt.GenerateToken(au.cfg.Server.Token.JWTKey, *clm)
	if err != nil {
		resonse.Code = fiber.StatusConflict
		resonse.Message = "Gagal membuat token"
		resonse.Message = "email atau password salah, silahkan coba lagi"
		resonse.Errors = err
		return resonse
	}

	redisKey := "authCust:token:" + strconv.Itoa(int(custBuild.ID))
	err = au.redis.Set(ctx, redisKey, token, 24*time.Hour).Err()
	if err != nil {
		resonse.Code = fiber.StatusConflict
		resonse.Message = "Gagal menyimpan token di Redis"
		resonse.Message = "email atau password salah, silahkan coba lagi"
		resonse.Errors = err
		return resonse
	}

	resBuild := new(responses.ResponseLogin)
	resBuild.ID = strconv.Itoa(int(custBuild.ID))
	resBuild.Token = token

	resonse.Data = &resBuild
	resonse.Code = fiber.StatusOK
	resonse.Message = "Pengguna berhasil masuk"
	return resonse
}

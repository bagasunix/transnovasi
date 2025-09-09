package usecases

import (
	"context"
	errs "errors"
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
	LoginUser(ctx context.Context, req *requests.Login) (resonses *responses.BaseResponse[*responses.ResponseLogin])
	LoginCustomer(ctx context.Context, req *requests.Login) (resonses *responses.BaseResponse[*responses.ResponseLogin])
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

func (au *authUsecase) LoginUser(ctx context.Context, req *requests.Login) (resonses *responses.BaseResponse[*responses.ResponseLogin]) {
	responseBuild := new(responses.BaseResponse[*responses.ResponseLogin])
	if req.Validate() != nil {
		responseBuild.Code = fiber.StatusBadRequest
		responseBuild.Message = "Validasi error"
		responseBuild.Errors = req.Validate()
		return responseBuild
	}
	// Check Login User
	checkUser := au.repo.GetUser().GetOneByParams(ctx, map[string]interface{}{"email": req.Email})
	if len(checkUser.Value.Email) == 0 || errs.Is(checkUser.Error, gorm.ErrRecordNotFound) {
		responseBuild.Code = fiber.StatusNotFound
		responseBuild.Message = "Email tidak ditemukan"
		responseBuild.Errors = errors.CustomError("email " + errors.ERR_NOT_FOUND)
		return responseBuild
	}

	if checkUser.Error != nil && !errs.Is(checkUser.Error, gorm.ErrRecordNotFound) {
		responseBuild.Code = fiber.StatusNotFound
		responseBuild.Message = checkUser.Error.Error()
		responseBuild.Errors = checkUser.Error
		return responseBuild
	}

	if !hash.ComparePasswords(checkUser.Value.Password, []byte(req.Password)) {
		responseBuild.Code = fiber.StatusNotFound
		responseBuild.Message = "username and password salah"
		responseBuild.Errors = errors.ErrInvalidAttributes("username and password")
		return responseBuild
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
		responseBuild.Code = fiber.StatusConflict
		responseBuild.Message = "Gagal membuat token"
		responseBuild.Errors = err
		return responseBuild
	}

	redisKey := "authUser:token:" + strconv.Itoa(int(userBuild.ID))
	err = au.redis.Set(ctx, redisKey, token, time.Hour).Err()
	if err != nil {
		responseBuild.Code = fiber.StatusConflict
		responseBuild.Message = "Gagal menyimpan token di Redis"
		responseBuild.Errors = err
		return responseBuild
	}

	resBuild := new(responses.ResponseLogin)
	resBuild.ID = strconv.Itoa(int(userBuild.ID))
	resBuild.Token = token

	responseBuild.Data = &resBuild
	responseBuild.Code = 200
	responseBuild.Message = "Pengguna berhasil masuk"
	return responseBuild
}

func (au *authUsecase) LoginCustomer(ctx context.Context, req *requests.Login) (resonses *responses.BaseResponse[*responses.ResponseLogin]) {
	responseBuild := new(responses.BaseResponse[*responses.ResponseLogin])
	if req.Validate() != nil {
		responseBuild.Code = fiber.StatusBadRequest
		responseBuild.Message = "Validasi error"
		responseBuild.Errors = req.Validate()
		return responseBuild
	}

	// Check Login Customer
	checkCust := au.repo.GetCustomer().GetOneByParams(ctx, map[string]interface{}{"email": req.Email})
	if len(checkCust.Value.Email) == 0 || errs.Is(checkCust.Error, gorm.ErrRecordNotFound) {
		responseBuild.Code = fiber.StatusNotFound
		responseBuild.Message = "Email tidak ditemukan"
		responseBuild.Errors = errors.CustomError("email " + errors.ERR_NOT_FOUND)
		return responseBuild
	}

	if checkCust.Error != nil && !errs.Is(checkCust.Error, gorm.ErrRecordNotFound) {
		responseBuild.Code = fiber.StatusNotFound
		responseBuild.Message = checkCust.Error.Error()
		responseBuild.Errors = checkCust.Error
		return responseBuild
	}

	if !hash.ComparePasswords(checkCust.Value.Password, []byte(req.Password)) {
		responseBuild.Code = fiber.StatusNotFound
		responseBuild.Message = "username and password salah"
		responseBuild.Errors = errors.ErrInvalidAttributes("username and password")
		return responseBuild
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
		responseBuild.Code = fiber.StatusConflict
		responseBuild.Message = "Gagal membuat token"
		responseBuild.Errors = err
		return responseBuild
	}

	redisKey := "authCust:token:" + strconv.Itoa(int(custBuild.ID))
	err = au.redis.Set(ctx, redisKey, token, 24*time.Hour).Err()
	if err != nil {
		responseBuild.Code = fiber.StatusConflict
		responseBuild.Message = "Gagal menyimpan token di Redis"
		responseBuild.Errors = err
		return responseBuild
	}

	resBuild := new(responses.ResponseLogin)
	resBuild.ID = strconv.Itoa(int(custBuild.ID))
	resBuild.Token = token

	responseBuild.Data = &resBuild
	responseBuild.Code = 200
	responseBuild.Message = "Pengguna berhasil masuk"
	return responseBuild
}

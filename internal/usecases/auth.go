package usecases

import (
	"context"
	errs "errors"
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
	LoginUser(ctx context.Context, request *requests.Login) (response responses.BaseResponse[*responses.ResponseLogin])
	Register(ctx context.Context, request *requests.User) (response responses.BaseResponse[*responses.UserResponse])
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
	userBuild.ID = strconv.Itoa(checkUser.Value.ID)
	userBuild.Name = checkUser.Value.Name
	userBuild.Sex = strconv.Itoa(checkUser.Value.Sex)
	userBuild.Email = checkUser.Value.Email
	userBuild.Role = checkUser.Value.Role
	userBuild.IsActive = strconv.Itoa(checkUser.Value.IsActive)

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

	redisKey := "auth_user:token:" + userBuild.ID
	err = au.redis.Set(ctx, redisKey, token, 24*time.Hour).Err()
	if err != nil {
		response.Code = fiber.StatusConflict
		response.Message = "email atau password salah, silahkan coba lagi"
		response.Errors = err.Error()
		return response
	}

	resBuild := new(responses.ResponseLogin)
	resBuild.ID = userBuild.ID
	resBuild.Token = token

	response.Data = &resBuild
	response.Code = fiber.StatusOK
	response.Message = "Pengguna berhasil masuk"
	return response
}

func (au *authUsecase) Register(ctx context.Context, request *requests.User) (response responses.BaseResponse[*responses.UserResponse]) {
	if request.Validate() != nil {
		response.Code = fiber.StatusBadRequest
		response.Message = "Validasi error"
		response.Errors = request.Validate().Error()
		return response
	}

	checkEmail := au.repo.GetUser().GetOneByParams(ctx, map[string]any{"email": request.Email})
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

	userBuild := new(domains.User)
	userBuild.Name = request.Name
	userBuild.Sex = request.Sex
	userBuild.Email = request.Email
	userBuild.Password = hash.HashAndSalt([]byte(request.Password))
	userBuild.IsActive = 1
	userBuild.Role = request.Role

	if err := au.repo.GetUser().Create(ctx, userBuild); err != nil {
		response.Code = fiber.StatusConflict
		response.Message = "Gagal membuat pengguna"
		response.Errors = err.Error()
		return response
	}

	resBuild := new(responses.UserResponse)
	resBuild.ID = strconv.Itoa(userBuild.ID)
	resBuild.Name = userBuild.Name
	resBuild.Sex = strconv.Itoa(userBuild.Sex)
	resBuild.Email = userBuild.Email
	resBuild.Role = userBuild.Role
	resBuild.IsActive = strconv.Itoa(userBuild.IsActive)

	response.Code = fiber.StatusCreated
	response.Message = "User berhasil dibuat"
	response.Data = &resBuild
	return response
}

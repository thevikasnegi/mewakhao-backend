package http

import (
	"ecom/internal/user/dto"
	"ecom/internal/user/service"
	"ecom/pkg/response"
	"ecom/pkg/utils"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gocommon/logger"
	"github.com/quangdangfit/gocommon/validation"
)

type UserController struct {
	srv       *service.UserService
	validator validation.Validation
}

func NewUserController(srv *service.UserService, validator validation.Validation) *UserController {
	return &UserController{
		srv:       srv,
		validator: validator,
	}
}

func (c *UserController) Login(ctx *gin.Context) {
	var req dto.LoginReq
	if err := ctx.ShouldBindJSON(&req); ctx.Request.Body == nil || err != nil {
		logger.Error("Failed to get body", err)
		response.Error(ctx, http.StatusBadRequest, err, "Invalid parameters")
		return
	}
	if err := c.validator.ValidateStruct(req); err != nil {
		logger.Error("Request body validation failed ", err)
		response.Error(ctx, http.StatusBadRequest, err, "Invalid parameters")
		return
	}
	user, accessToken, refreshToken, err := c.srv.Login(ctx, &req)
	if err != nil {
		logger.Error("Failed to login ", err)
		response.Error(ctx, http.StatusInternalServerError, err, "Something went wrong")
		return
	}

	var res dto.LoginRes
	utils.Copy(&res.User, &user)
	res.AccessToken = accessToken
	res.RefreshToken = refreshToken
	response.JSON(ctx, http.StatusOK, res)
}

func (c *UserController) Register(ctx *gin.Context) {
	var req dto.RegisterReq
	if err := ctx.ShouldBindJSON(&req); ctx.Request.Body == nil || err != nil {
		logger.Error("Failed to get body ", err)
		response.Error(ctx, http.StatusBadRequest, err, "Invalid parameters")
		return
	}
	if err := c.validator.ValidateStruct(req); err != nil {
		logger.Error("Request body validation failed ", err)
		response.Error(ctx, http.StatusBadRequest, err, "Invalid parameters")
		return
	}

	user, err := c.srv.Register(ctx, &req)
	if err != nil {
		logger.Error("Failed to register user", err)
		response.Error(ctx, http.StatusConflict, err, "Registeration Failed")
		return
	}
	var res dto.RegisterRes
	utils.Copy(&res.User, &user)
	response.JSON(ctx, http.StatusOK, res)
}

func (c *UserController) GetMe(ctx *gin.Context) {
	user, ok := ctx.Get("user")
	if !ok {
		response.Error(ctx, http.StatusBadRequest, errors.New("Invalid User"), "Invalid user")
		return
	}

	response.JSON(ctx, http.StatusOK, user)
}

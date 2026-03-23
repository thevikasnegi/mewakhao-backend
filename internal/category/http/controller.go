package http

import (
	"ecom/internal/category/dto"
	"ecom/internal/category/service"
	"ecom/pkg/response"
	"ecom/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gocommon/logger"
	"github.com/quangdangfit/gocommon/validation"
)

type CategoryController struct {
	srv       *service.CategoryService
	validator validation.Validation
}

func NewCategoryController(srv *service.CategoryService, validator validation.Validation) *CategoryController {
	return &CategoryController{
		srv:       srv,
		validator: validator,
	}
}

func (ctrl *CategoryController) List(ctx *gin.Context) {
	categories, err := ctrl.srv.List(ctx)
	if err != nil {
		logger.Error("Failed to list categories", err)
		response.Error(ctx, http.StatusInternalServerError, err, "Failed to list categories")
		return
	}

	var res []dto.CategoryRes
	utils.Copy(&res, &categories)
	response.JSON(ctx, http.StatusOK, res)
}

func (ctrl *CategoryController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	category, err := ctrl.srv.GetByID(ctx, id)
	if err != nil {
		logger.Error("Failed to get category", err)
		response.Error(ctx, http.StatusNotFound, err, "Category not found")
		return
	}

	var res dto.CategoryRes
	utils.Copy(&res, &category)
	response.JSON(ctx, http.StatusOK, res)
}

func (ctrl *CategoryController) Create(ctx *gin.Context) {
	var req dto.CreateCategoryReq
	if err := ctx.ShouldBindJSON(&req); ctx.Request.Body == nil || err != nil {
		logger.Error("Failed to get body", err)
		response.Error(ctx, http.StatusBadRequest, err, "Invalid parameters")
		return
	}
	if err := ctrl.validator.ValidateStruct(req); err != nil {
		logger.Error("Validation failed", err)
		response.Error(ctx, http.StatusBadRequest, err, "Invalid parameters")
		return
	}

	category, err := ctrl.srv.Create(ctx, &req)
	if err != nil {
		logger.Error("Failed to create category", err)
		response.Error(ctx, http.StatusConflict, err, "Failed to create category")
		return
	}

	var res dto.CategoryRes
	utils.Copy(&res, &category)
	response.JSON(ctx, http.StatusCreated, res)
}

func (ctrl *CategoryController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var req dto.UpdateCategoryReq
	if err := ctx.ShouldBindJSON(&req); ctx.Request.Body == nil || err != nil {
		logger.Error("Failed to get body", err)
		response.Error(ctx, http.StatusBadRequest, err, "Invalid parameters")
		return
	}

	category, err := ctrl.srv.Update(ctx, id, &req)
	if err != nil {
		logger.Error("Failed to update category", err)
		response.Error(ctx, http.StatusInternalServerError, err, "Failed to update category")
		return
	}

	var res dto.CategoryRes
	utils.Copy(&res, &category)
	response.JSON(ctx, http.StatusOK, res)
}

func (ctrl *CategoryController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := ctrl.srv.Delete(ctx, id); err != nil {
		logger.Error("Failed to delete category", err)
		response.Error(ctx, http.StatusInternalServerError, err, "Failed to delete category")
		return
	}

	response.JSON(ctx, http.StatusOK, nil)
}

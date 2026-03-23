package http

import (
	"ecom/internal/cart/dto"
	"ecom/internal/cart/service"
	"ecom/pkg/response"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gocommon/logger"
	"github.com/quangdangfit/gocommon/validation"
)

type CartController struct {
	srv       *service.CartService
	validator validation.Validation
}

func NewCartController(srv *service.CartService, validator validation.Validation) *CartController {
	return &CartController{
		srv:       srv,
		validator: validator,
	}
}

func getUserID(ctx *gin.Context) (string, error) {
	userID, exists := ctx.Get("userId")
	if !exists {
		return "", errors.New("user not authenticated")
	}
	id, ok := userID.(string)
	if !ok {
		return "", errors.New("invalid user ID")
	}
	return id, nil
}

func (ctrl *CartController) GetCart(ctx *gin.Context) {
	userID, err := getUserID(ctx)
	if err != nil {
		response.Error(ctx, http.StatusUnauthorized, err, "Unauthorized")
		return
	}

	cart, err := ctrl.srv.GetCart(ctx, userID)
	if err != nil {
		logger.Error("Failed to get cart", err)
		response.Error(ctx, http.StatusInternalServerError, err, "Failed to get cart")
		return
	}

	response.JSON(ctx, http.StatusOK, cart)
}

func (ctrl *CartController) AddItem(ctx *gin.Context) {
	userID, err := getUserID(ctx)
	if err != nil {
		response.Error(ctx, http.StatusUnauthorized, err, "Unauthorized")
		return
	}

	var req dto.AddItemReq
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

	cart, err := ctrl.srv.AddItem(ctx, userID, &req)
	if err != nil {
		logger.Error("Failed to add item to cart", err)
		response.Error(ctx, http.StatusInternalServerError, err, "Failed to add item")
		return
	}

	response.JSON(ctx, http.StatusOK, cart)
}

func (ctrl *CartController) UpdateItem(ctx *gin.Context) {
	userID, err := getUserID(ctx)
	if err != nil {
		response.Error(ctx, http.StatusUnauthorized, err, "Unauthorized")
		return
	}

	itemID := ctx.Param("itemId")
	var req dto.UpdateItemReq
	if err := ctx.ShouldBindJSON(&req); ctx.Request.Body == nil || err != nil {
		logger.Error("Failed to get body", err)
		response.Error(ctx, http.StatusBadRequest, err, "Invalid parameters")
		return
	}

	cart, err := ctrl.srv.UpdateItem(ctx, userID, itemID, &req)
	if err != nil {
		logger.Error("Failed to update cart item", err)
		response.Error(ctx, http.StatusInternalServerError, err, "Failed to update item")
		return
	}

	response.JSON(ctx, http.StatusOK, cart)
}

func (ctrl *CartController) RemoveItem(ctx *gin.Context) {
	userID, err := getUserID(ctx)
	if err != nil {
		response.Error(ctx, http.StatusUnauthorized, err, "Unauthorized")
		return
	}

	itemID := ctx.Param("itemId")
	cart, err := ctrl.srv.RemoveItem(ctx, userID, itemID)
	if err != nil {
		logger.Error("Failed to remove cart item", err)
		response.Error(ctx, http.StatusInternalServerError, err, "Failed to remove item")
		return
	}

	response.JSON(ctx, http.StatusOK, cart)
}

func (ctrl *CartController) ClearCart(ctx *gin.Context) {
	userID, err := getUserID(ctx)
	if err != nil {
		response.Error(ctx, http.StatusUnauthorized, err, "Unauthorized")
		return
	}

	if err := ctrl.srv.ClearCart(ctx, userID); err != nil {
		logger.Error("Failed to clear cart", err)
		response.Error(ctx, http.StatusInternalServerError, err, "Failed to clear cart")
		return
	}

	response.JSON(ctx, http.StatusOK, nil)
}

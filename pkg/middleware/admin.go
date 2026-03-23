package middleware

import (
	"ecom/pkg/response"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequireAdmin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		role, exists := ctx.Get("role")
		if !exists || role != "admin" {
			response.Error(ctx, http.StatusForbidden, errors.New("Forbidden"), "Admin access required")
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

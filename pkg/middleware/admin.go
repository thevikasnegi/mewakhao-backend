package middleware

import (
	"ecom/internal/user/entity"
	"ecom/pkg/response"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequireAdmin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		role, exists := ctx.Get("role")
		fmt.Printf("Role in require admin : %s\n", role)
		if !exists || role != entity.UserRoleAdmin {
			response.Error(ctx, http.StatusForbidden, errors.New("Forbidden"), "Admin access required")
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

package middleware

import (
	"ecom/internal/user/entity"
	"ecom/pkg/dbs"
	"ecom/pkg/jwt"
	"ecom/pkg/response"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func JWT(tokenType string, db *dbs.Database) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		if token == "" {
			response.Error(ctx, http.StatusUnauthorized, errors.New("Unauthorized"), "Unauthorized access")
			ctx.Abort()
			return
		}

		payload, err := jwt.ValidateToken(token)
		userId, ok := payload["id"].(string)
		if !ok || err != nil || payload == nil || payload["type"] != tokenType {
			response.Error(ctx, http.StatusUnauthorized, errors.New("Unauthorized"), "Unauthorized access")
			ctx.Abort()
			return
		}

		ctx.Set("userId", userId)
		ctx.Set("role", payload["role"])
		var user *entity.User

		db.FindById(ctx, userId, &user)
		ctx.Set("user", user)
		ctx.Next()
	}
}

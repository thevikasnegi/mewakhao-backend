package http

import (
	"ecom/internal/user/repository"
	"ecom/internal/user/service"
	"ecom/pkg/dbs"
	"ecom/pkg/jwt"
	"ecom/pkg/middleware"

	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gocommon/validation"
)

func Routes(r *gin.RouterGroup, db *dbs.Database, validator validation.Validation) {
	userRepo := repository.NewUserRepository(db)
	userSvc := service.NewUserService(validator, userRepo)
	userController := NewUserController(userSvc, validator)

	authMiddleware := middleware.JWT(jwt.AccessTokenType, db)
	authRoute := r.Group("/auth")
	{
		authRoute.POST("/register", userController.Register)
		authRoute.POST("/login", userController.Login)
		authRoute.GET("/me", authMiddleware, userController.GetMe)
	}
}

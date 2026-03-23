package http

import (
	"ecom/internal/category/repository"
	"ecom/internal/category/service"
	"ecom/pkg/dbs"
	"ecom/pkg/jwt"
	"ecom/pkg/middleware"

	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gocommon/validation"
)

func Routes(r *gin.RouterGroup, db *dbs.Database, validator validation.Validation) {
	categoryRepo := repository.NewCategoryRepository(db)
	categorySvc := service.NewCategoryService(categoryRepo)
	categoryController := NewCategoryController(categorySvc, validator)

	authMiddleware := middleware.JWT(jwt.AccessTokenType, db)
	adminMiddleware := middleware.RequireAdmin()

	categoryRoute := r.Group("/categories")
	{
		categoryRoute.GET("", categoryController.List)
		categoryRoute.GET("/:id", categoryController.GetByID)
		categoryRoute.POST("", authMiddleware, adminMiddleware, categoryController.Create)
		categoryRoute.PUT("/:id", authMiddleware, adminMiddleware, categoryController.Update)
		categoryRoute.DELETE("/:id", authMiddleware, adminMiddleware, categoryController.Delete)
	}
}

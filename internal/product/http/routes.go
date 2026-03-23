package http

import (
	"ecom/internal/product/repository"
	"ecom/internal/product/service"
	"ecom/pkg/dbs"
	"ecom/pkg/jwt"
	"ecom/pkg/middleware"

	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gocommon/validation"
)

func Routes(r *gin.RouterGroup, db *dbs.Database, validator validation.Validation) {
	productRepo := repository.NewProductRepository(db)
	productSvc := service.NewProductService(productRepo)
	productController := NewProductController(productSvc, validator)

	authMiddleware := middleware.JWT(jwt.AccessTokenType, db)
	adminMiddleware := middleware.RequireAdmin()

	productRoute := r.Group("/products")
	{
		productRoute.GET("", productController.List)
		productRoute.GET("/:id", productController.GetByID)
		productRoute.GET("/slug/:slug", productController.GetBySlug)
		productRoute.POST("", authMiddleware, adminMiddleware, productController.Create)
		productRoute.PUT("/:id", authMiddleware, adminMiddleware, productController.Update)
		productRoute.PUT("/:id/inventory", authMiddleware, adminMiddleware, productController.UpdateInventory)
		productRoute.DELETE("/:id", authMiddleware, adminMiddleware, productController.Delete)
	}
}

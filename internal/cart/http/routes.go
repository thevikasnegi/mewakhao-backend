package http

import (
	"ecom/internal/cart/repository"
	"ecom/internal/cart/service"
	"ecom/pkg/dbs"
	"ecom/pkg/jwt"
	"ecom/pkg/middleware"

	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gocommon/validation"
)

func Routes(r *gin.RouterGroup, db *dbs.Database, validator validation.Validation) {
	cartRepo := repository.NewCartRepository(db)
	cartSvc := service.NewCartService(cartRepo)
	cartController := NewCartController(cartSvc, validator)

	authMiddleware := middleware.JWT(jwt.AccessTokenType, db)

	cartRoute := r.Group("/cart")
	cartRoute.Use(authMiddleware)
	{
		cartRoute.GET("", cartController.GetCart)
		cartRoute.POST("/items", cartController.AddItem)
		cartRoute.PUT("/items/:itemId", cartController.UpdateItem)
		cartRoute.DELETE("/items/:itemId", cartController.RemoveItem)
		cartRoute.DELETE("", cartController.ClearCart)
	}
}

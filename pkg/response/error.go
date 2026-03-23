package response

import (
	"ecom/pkg/config"

	"github.com/gin-gonic/gin"
)

func Error(c *gin.Context, status int, err error, message string) {
	cfg := config.GetEnv()
	errorRes := map[string]interface{}{
		"message": message,
	}

	if cfg.Environment != config.ProductionEnv {
		errorRes["debug"] = err.Error()
	}

	c.JSON(status, Response{Error: errorRes})
}

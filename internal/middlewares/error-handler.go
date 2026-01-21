package middlewares

import (
	"shortner-url/infra/config"
	"shortner-url/internal"

	"github.com/gin-gonic/gin"
)

type Middlewares struct{}

func (m *Middlewares) ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		lastErr := c.Errors.Last()

		if lastErr == nil {
			return
		}

		if structured, ok := lastErr.Err.(*internal.APIError); ok {
			config.Logger().Error("Erro na requisição.", lastErr.Err, structured.Code)

			c.JSON(structured.Status, gin.H{
				"message": structured.Message,
				"code":    structured.Code,
			})

			return
		}

		config.Logger().Error("Erro na requisição.", lastErr.Err)

		c.JSON(500, gin.H{
			"message": "Erro interno do servidor.",
			"code":    000,
		})
	}
}

var Interceptors = Middlewares{}

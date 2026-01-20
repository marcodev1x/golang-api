package rest

import (
	"go-project/infra"
	"go-project/infra/config"
	"go-project/internal"
	"go-project/internal/repository/mysql"
	"go-project/internal/usecases"

	"github.com/gin-gonic/gin"
)

func UserRoutes() *[]internal.RouteHandler {
	rest := NewUserRest(usecases.NewUserUseCase(mysql.NewUserRespotory(infra.DomainDatabase)))

	config.Logger().Info(rest)

	return &[]internal.RouteHandler{
		{
			Path: "/users",
			Handler: func(c *gin.Context) {
				c.JSON(200, gin.H{
					"message": "Hello World",
				})
			},
			Method: internal.GET,
		},
	}
}

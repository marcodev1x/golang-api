package rest

import (
	"fmt"
	"shortner-url/infra"
	"shortner-url/internal"
	"shortner-url/internal/middlewares"
	"shortner-url/internal/repository/mysql"
	"shortner-url/internal/usecases"

	"github.com/gin-gonic/gin"
)

func UrlRoutes() *[]internal.RouteHandler {
	fmt.Println(infra.DomainDatabase)
	rest := NewUrlRest(usecases.NewUrlUseCase(mysql.NewUrlRepository(infra.DomainDatabase), usecases.NewRedisUsecase()))

	return &[]internal.RouteHandler{
		{
			Path:        "/url",
			Handler:     rest.FindUrlByHashedId,
			Method:      internal.GET,
			Middlewares: []gin.HandlerFunc{middlewares.Interceptors.ErrorHandler()},
		},
		{
			Path:        "/create-url",
			Handler:     rest.CreateUrl,
			Method:      internal.POST,
			Middlewares: []gin.HandlerFunc{middlewares.Interceptors.ErrorHandler()},
		},
	}
}

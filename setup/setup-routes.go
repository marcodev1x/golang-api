package setup

import (
	"shortner-url/internal"
	"shortner-url/internal/rest"

	"github.com/gin-gonic/gin"
)

func PrepareRoutes(server *gin.Engine) {
	internal.RouteDefiner(rest.UrlRoutes(), server)
}

package setup

import (
	"go-project/internal"
	"go-project/internal/rest"

	"github.com/gin-gonic/gin"
)

func PrepareRoutes(server *gin.Engine) {
	internal.RouteDefiner(rest.UserRoutes(), server)
}

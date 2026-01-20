package infra

import (
	"go-project/infra/config"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Bootstrap struct{}

func (b *Bootstrap) LoadEnv() *config.Env {
	env := config.LoadEnv()

	config.Logger().Info("Environment variables loaded")

	return env
}

func (b *Bootstrap) RunServer() *gin.Engine {
	router := gin.Default()

	config.Logger().Info("Starting server...")

	return router
}

func (b *Bootstrap) SetupDatabase(env *config.Env) *gorm.DB {
	database := &Database{}
	instance, err := database.Connect(env.DatabaseConfig)

	if err != nil {
		config.Logger().Error("Failed to connect to database", zap.Error(err))

		return nil
	}

	return instance
}

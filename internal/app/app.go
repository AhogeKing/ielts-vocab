package app

import (
	"ielts-vocab/internal/config"
	"ielts-vocab/internal/server"
	"ielts-vocab/internal/user"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type App struct {
	Router *gin.Engine
}

func New(db *gorm.DB, cfg *config.Config) *App {
	router := gin.Default()

	server.SetupRouter(
		router,
		user.NewModule(db, cfg))

	return &App{
		Router: router,
	}
}

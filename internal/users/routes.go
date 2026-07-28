package users

import (
	"Ielts-vocab/internal/auth"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB, tokens *auth.TokenManager) {
	usersRepo := &Repository{db: db}
	usersService := &Service{repo: usersRepo}
	usersHandler := &Handler{s: usersService, tokens: tokens}

	auth := rg.Group("/users")
	{
		auth.POST("/login", usersHandler.Login)
		auth.POST("/register", usersHandler.Register)
	}
}

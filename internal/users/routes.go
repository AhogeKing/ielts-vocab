package users

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	usersRepo := &Repository{db: db}
	usersService := &Service{repo: usersRepo}
	usersHandler := &Handler{usersService}

	auth := rg.Group("/users")
	{
		auth.POST("/login", usersHandler.Login)
		auth.POST("/register", usersHandler.Register)
	}
}

package user

import (
	"github.com/gin-gonic/gin"
)

func (m *Module) RegisterRoutes(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")
	{
		auth.POST("/login", m.handler.Login)
		auth.POST("/register", m.handler.Register)
	}
}

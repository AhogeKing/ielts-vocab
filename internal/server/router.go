package server

import (
	"Ielts-vocab/internal/users"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(r *gin.Engine, db *gorm.DB) {
	// /api
	api := r.Group("/api/v1")

	users.RegisterRoutes(api, db)
}

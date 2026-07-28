package server

import (
	"Ielts-vocab/internal/auth"
	"Ielts-vocab/internal/middleware"
	"Ielts-vocab/internal/users"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(r *gin.Engine, db *gorm.DB) error {
	tokens, err := auth.NewTokenManagerFromEnv()
	if err != nil {
		return err
	}

	// /api
	api := r.Group("/api/v1")
	api.Use(middleware.RequireJWT(
		tokens,
		"/api/v1/users/login",
		"/api/v1/users/register",
	))

	users.RegisterRoutes(api, db, tokens)
	return nil
}

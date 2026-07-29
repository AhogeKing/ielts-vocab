package server

import (
	"github.com/gin-gonic/gin"
)

type RouteRegister interface {
	RegisterRoutes(rg *gin.RouterGroup)
}

func SetupRouter(router *gin.Engine, modules ...RouteRegister) {
	// /api
	api := router.Group("/api/v1")

	for _, module := range modules {
		module.RegisterRoutes(api)
	}
}

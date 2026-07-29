package user

import (
	"ielts-vocab/internal/config"

	"gorm.io/gorm"
)

// Module represent the entire users
type Module struct {
	handler *Handler
}

func NewModule(db *gorm.DB, cfg *config.Config) *Module {
	repo := NewRepository(db)
	svc := NewService(cfg, repo)
	h := NewHandler(svc)

	return &Module{
		handler: h,
	}
}

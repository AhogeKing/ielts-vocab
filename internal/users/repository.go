package users

import (
	"context"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func (r *Repository) SelectByUsername(username string) (User, error) {
	ctx := context.Background()
	user, err := gorm.G[User](r.db).Where("username = ?", username).First(ctx)
	return user, err
}

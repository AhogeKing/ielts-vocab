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

func (r *Repository) Create(u User) error {
	ctx := context.Background()
	return gorm.G[User](r.db).Create(ctx, &u)
	// return r.db.WithContext(ctx).Create(&u).Error
}

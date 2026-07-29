package user

import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	FindByUsername(ctx context.Context, username string) (User, error)
	Create(ctx context.Context, user *User) error
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

type repository struct {
	db *gorm.DB
}

func (r *repository) FindByUsername(ctx context.Context, username string) (User, error) {
	user, err := gorm.G[User](r.db).Where("username = ?", username).First(ctx)
	return user, err
}

func (r *repository) Create(ctx context.Context, u *User) error {
	return gorm.G[User](r.db).Create(ctx, u)
	// return r.db.WithContext(ctx).Create(&u).Error
}

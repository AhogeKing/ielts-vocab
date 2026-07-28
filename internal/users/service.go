package users

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type Service struct {
	repo *Repository
}

func (s *Service) Login(l LoginRequest) (User, error) {
	user, err := s.repo.SelectByUsername(l.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, &ErrInvalidCredentials{}
		}
		return user, err
	}

	if err := comparePassword(user.PasswordHash, l.Password); err != nil {
		return user, &ErrInvalidCredentials{}
	}

	if user.IsActive == false {
		return user, &ErrAccountDisabled{}
	}

	return user, nil
}

func (s *Service) Register(r RegisterRequest) error {
	_, err := s.repo.SelectByUsername(r.Username)

	switch {
	case err == nil:
		// 查询成功，说明数据库已有该用户名，不能重复注册
		return &ErrUsernameAlreadyExist{}

	case errors.Is(err, gorm.ErrRecordNotFound):
		passwordHash, err := hashPassword(r.Password)
		if err != nil {
			return fmt.Errorf("hash password: %w", err)
		}

		user := User{
			Username:     r.Username,
			PasswordHash: passwordHash,
			IsActive:     true,
		}
		if r.Email != "" {
			email := r.Email
			user.Email = &email
		}

		if err := s.repo.Create(user); err != nil {
			return fmt.Errorf("create user: %w", err)
		}

		return nil

	default:
		return fmt.Errorf("check username: %w", err)
	}
}

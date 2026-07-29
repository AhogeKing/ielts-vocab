package user

import (
	"context"
	"errors"
	"fmt"
	"ielts-vocab/internal/config"

	"gorm.io/gorm"
)

type Service interface {
	Register(ctx context.Context, req RegisterRequest) error
	Login(ctx context.Context, req LoginRequest) (LoginResponse, error)
}

type service struct {
	cfg  *config.Config
	repo Repository
}

func NewService(cfg *config.Config, repo Repository) Service {
	return &service{
		cfg:  cfg,
		repo: repo,
	}
}

func (s *service) Register(ctx context.Context, req RegisterRequest) error {
	_, err := s.repo.FindByUsername(ctx, req.Username)

	switch {
	case err == nil:
		// check whether user already exists
		// 查询成功，说明数据库已有该用户名，不能重复注册
		return ErrUsernameAlreadyExist

	case errors.Is(err, gorm.ErrRecordNotFound):
		// hash password
		passwordHash, err := hashPassword(req.Password)
		if err != nil {
			return fmt.Errorf("hash password: %w", err)
		}

		user := &User{
			Username:     req.Username,
			PasswordHash: passwordHash,
			IsActive:     true,
		}
		if req.Email != "" {
			email := req.Email
			user.Email = &email
		}

		// create user
		if err := s.repo.Create(ctx, user); err != nil {
			return fmt.Errorf("create user: %w", err)
		}
		return nil

	default:
		return fmt.Errorf("check username: %w", err)
	}
}

func (s *service) Login(ctx context.Context, req LoginRequest) (LoginResponse, error) {
	// check user exists
	user, err := s.repo.FindByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return LoginResponse{}, ErrInvalidCredentials
		}
		return LoginResponse{}, err
	}

	if err = comparePassword(user.PasswordHash, req.Password); err != nil {
		return LoginResponse{}, ErrInvalidCredentials
	}

	// generate access token

	// get refresh token if exists

	// generate & store refresh token

	// return
}

//func (s *Service) Login(l LoginRequest) (User, error) {
//	user, err := s.repo.SelectByUsername(l.Username)
//	if err != nil {
//		if errors.Is(err, gorm.ErrRecordNotFound) {
//			return user, ErrInvalidCredentials
//		}
//		return user, err
//	}
//
//	if err := comparePassword(user.PasswordHash, l.Password); err != nil {
//		return user, ErrInvalidCredentials
//	}
//
//	if user.IsActive == false {
//		return user, ErrAccountDisabled
//	}
//
//	return user, nil
//}

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
			fmt.Printf("user not found with username %s", l.Username)
			return user, &ErrInvalidCredentials{}
		}
		return user, err
	}

	if user.PasswordHash != l.Password {
		return user, &ErrInvalidCredentials{}
	}

	if user.IsActive == false {
		return user, &ErrAccountDisabled{}
	}

	return user, nil
}

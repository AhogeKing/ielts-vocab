package users

import "errors"

var (
	ErrInvalidCredentials   = errors.New("username or password is invalid")
	ErrAccountDisabled      = errors.New("the account is disabled")
	ErrUsernameAlreadyExist = errors.New("username already exists")
)

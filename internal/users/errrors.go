package users

type ErrInvalidCredentials struct {
}

func (e *ErrInvalidCredentials) Error() string {
	return "username or password is invalid"
}

type ErrAccountDisabled struct{}

func (e *ErrAccountDisabled) Error() string {
	return "the account is disabled"
}

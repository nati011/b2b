package auth

import "errors"

// system errors
var (
	ErrUsernameTaken = errors.New("invalid username")
	ErrEmailTaken    = errors.New("invalid email")
	ErrFailedToLogin = errors.New("invalid email or password")
)

type LoginAuthResonse struct {
}

type CreateClientAuthResonse struct {
}

type AuthProvider interface {
	CreateNewClient(firstName string, lastName string, email string, username string, password string) (CreateClientAuthResonse, error)
	ClientLogin(email, password string) (LoginAuthResonse, error)
}

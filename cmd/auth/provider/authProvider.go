package auth

import "errors"

// system errors
var (
	ErrSysUsernameTaken = errors.New("invalid username")
	ErrSysEmailTaken    = errors.New("invalid email")
	ErrSysFailedToLogin = errors.New("invalid email or password")
)

type LoginAuthResonse struct {
	JWT JWT
}

type JWT struct {
	AccessToken      string
	IDToken          string
	ExpiresIn        int
	RefreshExpiresIn int
	RefreshToken     string
	TokenType        string
	NotBeforePolicy  int
	SessionState     string
	Scope            string
}

type CreateClientAuthResonse struct {
	Username string
}

type AuthProvider interface {
	CreateNewClient(firstName string, lastName string, email string, username string, password string) (CreateClientAuthResonse, error)
	ClientLogin(email, password string) (LoginAuthResonse, error)
}

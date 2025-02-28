package provider

import "errors"

// system errors
var (
	ErrSysUsernameTaken = errors.New("invalid username")
	ErrSysEmailTaken    = errors.New("invalid email")
	ErrSysFailedToLogin = errors.New("invalid email or password")
	ErrSysTokenExpired  = errors.New("invalid or expired token")
	ErrSysUnknown       = errors.New("unknown error")
)

type LoginAuthResponse struct {
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

type CreateClientAuthResponse struct {
	Username string
}

type AuthProvider interface {
	CreateNewClient(firstName string, lastName string, email string, username string, password string) (CreateClientAuthResponse, error)
	ClientLogin(email, password string) (LoginAuthResponse, error)
	RefreshToken(token string) (LoginAuthResponse, error)
}

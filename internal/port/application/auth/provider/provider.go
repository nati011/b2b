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

type RegisterUserRequest struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmed_password"`
	FullName        string `json:"full_name"`
	Username        string `json:"username"`
}

type RegisterUserResponse struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type LoginAuthResponse struct {
	JWT     JWT
	Message string `json:"message"`
}

type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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

type Provider interface {
	CreateNewClient(firstName string, lastName string, email string, username string, password string) (CreateClientAuthResponse, error)
	ClientLogin(email, password string) (LoginAuthResponse, error)
}

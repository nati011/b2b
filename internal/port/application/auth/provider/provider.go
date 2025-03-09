package provider

import (
	"context"
	"errors"
	"time"
)

// system errors
var (
	ErrSysUsernameTaken = errors.New("invalid username")
	ErrSysEmailTaken    = errors.New("invalid email")
	ErrSysFailedToLogin = errors.New("invalid email or password")
	ErrSysTokenExpired  = errors.New("invalid or expired token")
	ErrSysUnknown       = errors.New("unknown error")
)

type RegisterUserRequest struct {
	Email           string    `json:"email"`
	Password        string    `json:"password"`
	ConfirmPassword string    `json:"confirmed_password"`
	BirthDate       time.Time `json:"birth_date"`
	PhoneNumber     string    `json:"phone_number"`
	ExternalId      string    `json:"external_id"`
	FirstName       string    `json:"first_name"`
	LastName        string    `json:"last_name"`
	Username        string    `json:"username"`
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
	CreateNewClient(ctx context.Context, req RegisterUserRequest) (CreateClientAuthResponse, error)
	ClientLogin(ctx context.Context, req LoginUserRequest) (LoginAuthResponse, error)
	RefreshToken(ctx context.Context, req RefreshTokenRequest) (LoginAuthResponse, error)
}

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
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 036df0cf (+ remove password confirmation)
	Email       string
	Password    string
	BirthDate   time.Time
	PhoneNumber string
	ExternalId  string
	FirstName   string
	LastName    string
	Username    string
<<<<<<< HEAD
}

type RegisterUserResponse struct {
	Id       string
=======
	Email           string
	Password        string
	ConfirmPassword string
	BirthDate       time.Time
	PhoneNumber     string
	ExternalId      string
	FirstName       string
	LastName        string
	Username        string
=======
>>>>>>> 036df0cf (+ remove password confirmation)
}

type RegisterUserResponse struct {
>>>>>>> 1734bfa2 (resolve conflict)
	Username string
}

type RefreshTokenRequest struct {
	RefreshToken string
}

type LoginAuthResponse struct {
	JWT JWT
}

type LoginUserRequest struct {
	Email    string
	Password string
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

type Provider interface {
	CreateNewClient(ctx context.Context, req RegisterUserRequest) (RegisterUserResponse, error)
	ClientLogin(ctx context.Context, req LoginUserRequest) (LoginAuthResponse, error)
	RefreshToken(ctx context.Context, req RefreshTokenRequest) (LoginAuthResponse, error)
	DeleteClient(ctx context.Context, userId string) error
}

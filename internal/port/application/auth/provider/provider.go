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
	Email       string
	Password    string
	BirthDate   time.Time
	PhoneNumber string
	ExternalId  string
	FirstName   string
	LastName    string
	Username    string
}

type RegisterUserResponse struct {
	Id       string
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

type RetrospectionResult struct {
	Exp      *int
	Nbf      *int
	Iat      *int
	Active   *bool
	AuthTime *int
	Jti      *string
	Type     *string
}

type Claims struct {
	ExpirationTime time.Time
	IssuedAt       time.Time
	NotBefore      time.Time
	Issuer         string
	Subject        string
	Audience       string
	Email          string
}

type DecodedResult struct {
	Raw       string
	Header    map[string]interface{}
	Signature []byte
	Valid     bool
	Claims    Claims
}

type AssignRoleRequest struct {
}

type Provider interface {
	CreateNewClient(ctx context.Context, req *RegisterUserRequest) (RegisterUserResponse, error)
	ClientLogin(ctx context.Context, req *LoginUserRequest) (LoginAuthResponse, error)
	RefreshToken(ctx context.Context, req *RefreshTokenRequest) (LoginAuthResponse, error)
	DeleteClient(ctx context.Context, userId string) error
	ResetPassword(ctx context.Context, userId, new_password string) error
	RetrospectToken(ctx context.Context, token string) (RetrospectionResult, error)
	DecodeToken(ctx context.Context, token string) (DecodedResult, error)
	ClientLogout(ctx context.Context, req RefreshTokenRequest) error
	AssignRole(ctx context.Context, userId string, roleId int) error
}

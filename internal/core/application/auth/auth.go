package auth

import (
	"context"
	"time"
)

type RegisterUserRequest struct {
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	BirthDate   time.Time `json:"birth_date"`
	PhoneNumber string    `json:"phone_number"`
	ExternalId  string    `json:"external_id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Username    string    `json:"username"`
}

type RegisterUserWithoutPasswordRequest struct {
	Email       string    `json:"email"`
	BirthDate   time.Time `json:"birth_date"`
	PhoneNumber string    `json:"phone_number"`
	ExternalId  string    `json:"external_id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Username    string    `json:"username"`
}

type RegisterUserResponse struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type LoginAuthResponse struct {
	JWT JWT `json:"jwt"`
}

type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ResetCredentialsRequest struct {
	NewPassword string `json:"password"`
	ResetToken  string `json:"resetToken"`
}

type InitClientCredentialsResetRequest struct {
	UserId string `json:"user_id"`
	Email  string `json:"email"`
}

type JWT struct {
	AccessToken      string `json:"access_token"`
	IDToken          string `json:"id_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	NotBeforePolicy  int    `json:"not_before_policy"`
	SessionState     string `json:"session_state"`
	Scope            string `json:"scope"`
}

type ResourceAccess struct {
	Roles []string `json:"roles"`
}

type Provider interface {
	CreateNewClientWithPassword(ctx context.Context, req RegisterUserRequest) (RegisterUserResponse, error)
	CreateNewClientWithOutPassword(ctx context.Context, req RegisterUserWithoutPasswordRequest) (RegisterUserResponse, error)
	ClientLogin(ctx context.Context, req LoginUserRequest) (LoginAuthResponse, error)
	RefreshToken(ctx context.Context, req RefreshTokenRequest) (LoginAuthResponse, error)
	DeleteClient(ctx context.Context, userId string) error
	ResetClientCredentials(ctx context.Context, req ResetCredentialsRequest) error
	IsAuthorizedForResource(ctx context.Context, role string) (bool, error)
	GetUserAuthorization(ctx context.Context) (ResourceAccess, error)
}

type Auth struct {
}

func NewAuthService() Provider {
	return &Auth{}
}

func (a *Auth) CreateNewClientWithPassword(ctx context.Context, req RegisterUserRequest) (RegisterUserResponse, error)
func (a *Auth) CreateNewClientWithOutPassword(ctx context.Context, req RegisterUserWithoutPasswordRequest) (RegisterUserResponse, error)
func (a *Auth) ClientLogin(ctx context.Context, req LoginUserRequest) (LoginAuthResponse, error)
func (a *Auth) RefreshToken(ctx context.Context, req RefreshTokenRequest) (LoginAuthResponse, error)
func (a *Auth) DeleteClient(ctx context.Context, userId string) error
func (a *Auth) ResetClientCredentials(ctx context.Context, req ResetCredentialsRequest) error
func (a *Auth) IsAuthorizedForResource(ctx context.Context, role string) (bool, error)
func (a *Auth) GetUserAuthorization(ctx context.Context) (ResourceAccess, error)

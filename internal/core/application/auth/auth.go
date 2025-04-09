package auth

import (
	"context"
	"errors"
	"time"

	port "b2b.nati011.github.com/internal/port/application/auth/provider"
)

var (
	ErrUsernameTaken        = errors.New("oopsy, username already taken")
	ErrEmailTaken           = errors.New("oopsy, email already taken")
	ErrFailedToLogin        = errors.New("oopsy, email or password incorrect")
	ErrUnknown              = errors.New("oopsy, unknown error has occured")
	ErrEmailNotSupplied     = errors.New("oopsy, email mandatory")
	ErrUsernameNotSupplied  = errors.New("oopsy, username mandatory")
	ErrInvalidEmail         = errors.New("oopsy, Invalid Email")
	ErrPasswordNotSupplied  = errors.New("oopsy, password mandatory")
	ErrFirstNameNotSupplied = errors.New("oopsy, First Name mandatory")
	ErrLastNameNotSupplied  = errors.New("oppsy, Last Name mandatory")
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

type RegisterUserResponse struct {
	Username string `json:"username"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type LoginAuthResponse struct {
	JWT JWT
}

type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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

type Provider interface {
	CreateNewClient(ctx context.Context, req RegisterUserRequest) (RegisterUserResponse, error)
	ClientLogin(ctx context.Context, req LoginUserRequest) (LoginAuthResponse, error)
	RefreshToken(ctx context.Context, req RefreshTokenRequest) (LoginAuthResponse, error)
}

type AuthService struct {
	authProvider port.Provider
}

func NewAuthService(ap port.Provider) Provider {
	return &AuthService{authProvider: ap}

}

func (a *AuthService) CreateNewClient(ctx context.Context, req RegisterUserRequest) (RegisterUserResponse, error) {
	err := validateName(req.FirstName, req.LastName)
	if err != nil {
		return RegisterUserResponse{}, err
	}

	err = validatePasswords(req.Password)
	if err != nil {
		return RegisterUserResponse{}, err
	}
	err = validateEmail(req.Email)
	if err != nil {
		return RegisterUserResponse{}, err
	}
	err = validateUsername(req.Username)
	if err != nil {
		return RegisterUserResponse{}, err
	}
	resp, err := a.authProvider.CreateNewClient(ctx, port.RegisterUserRequest(req))
	if err != nil {
		switch err {
		case port.ErrSysUsernameTaken:
			return RegisterUserResponse{}, ErrUsernameTaken
		case port.ErrSysEmailTaken:
			return RegisterUserResponse{}, ErrEmailTaken
		default:
			return RegisterUserResponse{}, ErrUnknown
		}
	}
	return RegisterUserResponse{
		Username: resp.Username,
	}, nil
}

func (a *AuthService) ClientLogin(ctx context.Context, rq LoginUserRequest) (LoginAuthResponse, error) {
	resp, err := a.authProvider.ClientLogin(ctx, port.LoginUserRequest(rq))
	if err != nil {
		switch err {
		case port.ErrSysFailedToLogin:
			return LoginAuthResponse{}, ErrFailedToLogin
		default:
			return LoginAuthResponse{}, ErrUnknown
		}
	}
	return LoginAuthResponse{
		JWT: JWT(resp.JWT),
	}, nil
}
func (a *AuthService) RefreshToken(ctx context.Context, req RefreshTokenRequest) (LoginAuthResponse, error) {
	resp, err := a.authProvider.RefreshToken(ctx, port.RefreshTokenRequest(req))
	if err != nil {
		switch err {
		case port.ErrSysFailedToLogin:
			return LoginAuthResponse{}, ErrFailedToLogin
		default:
			return LoginAuthResponse{}, ErrUnknown
		}
	}
	return LoginAuthResponse{
		JWT: JWT(resp.JWT),
	}, nil
}

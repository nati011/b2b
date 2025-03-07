package auth

import (
	"context"
	"errors"

	port "b2b.nati011.github.com/internal/port/application/auth/provider"
)

var (
	ErrUsernameTaken                   = errors.New("oopsy, username already taken")
	ErrEmailTaken                      = errors.New("oopsy, email already taken")
	ErrFailedToLogin                   = errors.New("oopsy, email or password incorrect")
	ErrUnknown                         = errors.New("oopsy, unknown error has occured")
	ErrEmailNotSupplied                = errors.New("oopsy, email mandatory")
	ErrPasswordNotSupplied             = errors.New("oopsy, password mandatory")
	ErrConfirmationPasswordNotSupplied = errors.New("oopsy, confirmation password mandatory")
	ErrFullNameNotSupplied             = errors.New("oopsy, fullname mandatory")
	ErrPasswordsDontMatch              = errors.New("oopsy, passwords dont match")
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
}

type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUserResonse struct {
	JWT JWT
}

type JWT struct {
	AccessToken      string `json:"access_token"`
	IDToken          string `json:"id_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	NotBeforePolicy  int    `json:"not-before-policy"`
	SessionState     string `json:"session_state"`
	Scope            string `json:"scope"`
}

type Provider interface {
	CreateClient(context.Context, RegisterUserRequest) (RegisterUserResponse, error)
	LoginClient(context.Context, LoginUserRequest) (LoginUserResonse, error)
}

type AuthService struct {
	authProvider port.Provider
}

func NewAuthService(ap port.Provider) Provider {
	return &AuthService{authProvider: ap}
}

func (a AuthService) CreateClient(ctx context.Context, req RegisterUserRequest) (RegisterUserResponse, error) {
	err := validatePasswords(req.Password, req.ConfirmPassword)
	if err != nil {
		return RegisterUserResponse{}, err
	}
	err = validateEmail(req.Email)
	if err != nil {
		return RegisterUserResponse{}, err
	}
	err = validateFullName(req.FullName)
	if err != nil {
		return RegisterUserResponse{}, err
	}

	resp, err := a.authProvider.CreateNewClient(req.FullName, req.FullName, req.Email, req.Username, req.Password)
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

func (a *AuthService) LoginClient(ctx context.Context, rq LoginUserRequest) (LoginUserResonse, error) {
	resp, err := a.authProvider.ClientLogin(rq.Email, rq.Password)
	if err != nil {
		switch err {
		case port.ErrSysFailedToLogin:
			return LoginUserResonse{}, ErrFailedToLogin
		default:
			return LoginUserResonse{}, ErrUnknown
		}
	}
	return LoginUserResonse{
		JWT: JWT(resp.JWT),
	}, nil
}

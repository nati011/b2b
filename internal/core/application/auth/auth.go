package auth

import (
	"context"
	"errors"
	"log"

	port "b2b.nati011.github.com/internal/port/application/auth/provider"
)

var (
	ErrUsernameTaken                   = errors.New("oopsy, username already taken")
	ErrEmailTaken                      = errors.New("oopsy, email already taken")
	ErrFailedToLogin                   = errors.New("oopsy, email or password incorrect")
	ErrUnknown                         = errors.New("oopsy, unknown error has occured")
	ErrEmailNotSupplied                = errors.New("oopsy, email mandatory")
	ErrInvalidEmail                    = errors.New("oopsy, Invalid Email")
	ErrPasswordNotSupplied             = errors.New("oopsy, password mandatory")
	ErrConfirmationPasswordNotSupplied = errors.New("oopsy, confirmation password mandatory")
	ErrFirstNameNotSupplied            = errors.New("oopsy, First Name mandatory")
	ErrLastNameNotSupplied             = errors.New("oppsy, Last Name mandatory")
	ErrPasswordsDontMatch              = errors.New("oopsy, passwords dont match")

	SUCCESS_MESSAGE = "Ahoy!"
)

type Provider interface {
	CreateNewClient(ctx context.Context, req port.RegisterUserRequest) (port.RegisterUserResponse, error)
	ClientLogin(ctx context.Context, req port.LoginUserRequest) (port.LoginAuthResponse, error)
	RefreshToken(ctx context.Context, req port.RefreshTokenRequest) (port.LoginAuthResponse, error)
}

type AuthService struct {
	authProvider Provider
}

func NewAuthService(ap port.Provider) port.Provider {
	return &AuthService{authProvider: ap}

}

func (a *AuthService) CreateNewClient(ctx context.Context, req port.RegisterUserRequest) (port.RegisterUserResponse, error) {
	err := validateName(req.FirstName, req.LastName)
	if err != nil {
		return port.RegisterUserResponse{}, err
	}

	err = validatePasswords(req.Password, req.ConfirmPassword)
	if err != nil {
		return port.RegisterUserResponse{}, err
	}
	err = validateEmail(req.Email)
	if err != nil {
		return port.RegisterUserResponse{}, err
	}

	resp, err := a.authProvider.CreateNewClient(ctx, req)
	if err != nil {
		switch err {
		case port.ErrSysUsernameTaken:
			return port.RegisterUserResponse{}, ErrUsernameTaken
		case port.ErrSysEmailTaken:
			return port.RegisterUserResponse{}, ErrEmailTaken
		default:
			return port.RegisterUserResponse{}, ErrUnknown
		}
	}
	return port.RegisterUserResponse{
		Username: resp.Username,
	}, nil
}

func (a *AuthService) ClientLogin(ctx context.Context, rq port.LoginUserRequest) (port.LoginAuthResponse, error) {
	resp, err := a.authProvider.ClientLogin(ctx, rq)
	log.Printf(resp.Message)
	if err != nil {
		switch err {
		case port.ErrSysFailedToLogin:
			return port.LoginAuthResponse{}, ErrFailedToLogin
		default:
			return port.LoginAuthResponse{}, ErrUnknown
		}
	}
	return port.LoginAuthResponse{
		JWT: port.JWT(resp.JWT),
	}, nil
}
func (a *AuthService) RefreshToken(ctx context.Context, req port.RefreshTokenRequest) (port.LoginAuthResponse, error) {
	resp, err := a.authProvider.RefreshToken(ctx, req)
	if err != nil {
		switch err {
		case port.ErrSysFailedToLogin:
			return port.LoginAuthResponse{
				Message: ErrFailedToLogin.Error(),
			}, ErrFailedToLogin
		default:
			return port.LoginAuthResponse{}, ErrUnknown
		}
	}
	return port.LoginAuthResponse{
		JWT:     port.JWT(resp.JWT),
		Message: SUCCESS_MESSAGE,
	}, nil
}

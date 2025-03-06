package auth

import (
	"errors"

	port "b2b.nati011.github.com/internal/port/application/auth/provider"
)

// user readable errors
var (
	ErrUsernameTaken = errors.New("oopsy, username already taken")
	ErrEmailTaken    = errors.New("oopsy, email already taken")
	ErrFailedToLogin = errors.New("oopsy, email or password incorrect")
	ErrUnknown       = errors.New("oopsy, unknown error has occured")
)

const (
	SUCCESS_MESSAGE = "Ahoy!"
)

type Authorizer interface {
	CreateClient(port.RegisterUserRequest) (port.RegisterUserResponse, error)
	LoginClient(port.LoginUserRequest) (port.LoginAuthResponse, error)
}

type AuthService struct {
	authProvider port.AuthProvider
}

func NewAuthService(ap port.AuthProvider) *AuthService {
	return &AuthService{authProvider: ap}
}

func (a AuthService) CreateClient(rq port.RegisterUserRequest) (port.RegisterUserResponse, error) {
	resp, err := a.authProvider.CreateNewClient(rq.FullName, rq.FullName, rq.Email, rq.Username, rq.Password)
	if err != nil {
		switch err {
		case port.ErrSysUsernameTaken:
			return port.RegisterUserResponse{
				Message: ErrUsernameTaken.Error(),
			}, ErrUsernameTaken
		case port.ErrSysEmailTaken:
			return port.RegisterUserResponse{
				Message: ErrEmailTaken.Error(),
			}, ErrEmailTaken
		default:
			return port.RegisterUserResponse{}, ErrUnknown
		}
	}
	return port.RegisterUserResponse{
		Username: resp.Username,
		Message:  SUCCESS_MESSAGE,
	}, nil
}

func (a *AuthService) LoginClient(rq port.LoginUserRequest) (port.LoginAuthResponse, error) {
	resp, err := a.authProvider.ClientLogin(rq.Email, rq.Password)
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

func (a *AuthService) RefreshToken(req port.RefreshTokenRequest) (port.LoginAuthResponse, error) {
	resp, err := a.authProvider.RefreshToken(req.RefreshToken)
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

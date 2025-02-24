package auth

import (
	"errors"

	authDTO "b2b.nati011.github.com/internal/core/application/util/auth/model/dto"
	provider "b2b.nati011.github.com/internal/core/application/util/auth/provider"
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
	CreateClient(authDTO.RegisterUserRequest) (authDTO.RegisterUserResponse, error)
	LoginClient(authDTO.LoginUserRequest) (authDTO.LoginUserResonse, error)
}

type AuthService struct {
	authProvider provider.AuthProvider
}

func NewAuthService(ap provider.AuthProvider) *AuthService {
	return &AuthService{authProvider: ap}
}

func (a AuthService) CreateClient(rq authDTO.RegisterUserRequest) (authDTO.RegisterUserResponse, error) {
	resp, err := a.authProvider.CreateNewClient(rq.FullName, rq.FullName, rq.Email, rq.Username, rq.Password)
	if err != nil {
		switch err {
		case provider.ErrSysUsernameTaken:
			return authDTO.RegisterUserResponse{
				Message: ErrUsernameTaken.Error(),
			}, ErrUsernameTaken
		case provider.ErrSysEmailTaken:
			return authDTO.RegisterUserResponse{
				Message: ErrEmailTaken.Error(),
			}, ErrEmailTaken
		default:
			return authDTO.RegisterUserResponse{}, ErrUnknown
		}
	}
	return authDTO.RegisterUserResponse{
		Username: resp.Username,
		Message:  SUCCESS_MESSAGE,
	}, nil
}

func (a *AuthService) LoginClient(rq authDTO.LoginUserRequest) (authDTO.LoginUserResonse, error) {
	resp, err := a.authProvider.ClientLogin(rq.Email, rq.Password)
	if err != nil {
		switch err {
		case provider.ErrSysFailedToLogin:
			return authDTO.LoginUserResonse{
				Message: ErrFailedToLogin.Error(),
			}, ErrFailedToLogin
		default:
			return authDTO.LoginUserResonse{}, ErrUnknown
		}
	}
	return authDTO.LoginUserResonse{
		JWT:     authDTO.JWT(resp.JWT),
		Message: SUCCESS_MESSAGE,
	}, nil
}

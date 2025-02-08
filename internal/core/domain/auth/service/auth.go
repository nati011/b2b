package auth

import (
	"errors"

	auth "b2b.nati011.github.com/internal/core/domain/auth/model/dto"
	provider "b2b.nati011.github.com/internal/core/domain/auth/provider"
)

// user readable errors
var (
	ErrUsernameTaken = errors.New("oopsy, username already taken")
	ErrEmailTaken    = errors.New("oopsy, email is already taken")
	ErrFailedToLogin = errors.New("oopsy, email or password incorrect")
	ErrUnknown       = errors.New("oopsy, unknown error has occured")
)

const (
	SUCCESS_MESSAGE = "Ahoy!"
)

type Authorizer interface {
	CreateClient(auth.RegisterUserRequest) (auth.RegisterUserResponse, error)
	LoginClient(auth.LoginUserRequest) (auth.LoginUserResonse, error)
}

type AuthService struct {
	authProvider provider.AuthProvider
}

func NewAuthService(ap provider.AuthProvider) *AuthService {
	return &AuthService{authProvider: ap}
}

func (a AuthService) CreateClient(rq auth.RegisterUserRequest) (auth.RegisterUserResponse, error) {
	resp, err := a.authProvider.CreateNewClient(rq.FullName, rq.FullName, rq.Email, rq.Username, rq.Password)
	if err != nil {
		switch err {
		case provider.ErrSysUsernameTaken:
			return auth.RegisterUserResponse{
				Username: "",
				Message:  ErrUsernameTaken.Error(),
			}, ErrUsernameTaken
		case provider.ErrSysEmailTaken:
			return auth.RegisterUserResponse{
				Username: "",
				Message:  ErrEmailTaken.Error(),
			}, ErrEmailTaken
		default:
			return auth.RegisterUserResponse{
				Username: "",
				Message:  ErrUnknown.Error(),
			}, ErrUnknown
		}
	}
	return auth.RegisterUserResponse{
		Username: resp.Username,
		Message:  SUCCESS_MESSAGE,
	}, nil
}

func (a *AuthService) LoginClient(rq auth.LoginUserRequest) (auth.LoginUserResonse, error) {
	resp, err := a.authProvider.ClientLogin(rq.Email, rq.Password)
	if err != nil {
		switch err {
		case provider.ErrSysFailedToLogin:
			return auth.LoginUserResonse{}, ErrFailedToLogin
		default:
			return auth.LoginUserResonse{}, ErrUnknown
		}
	}
	return auth.LoginUserResonse{
		JWT:     auth.JWT(resp.JWT),
		Message: SUCCESS_MESSAGE,
	}, nil
}

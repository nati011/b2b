package auth

import (
	"errors"
	"log"

	auth "b2b.nati011.github.com/cmd/auth/model/dto"
	provider "b2b.nati011.github.com/cmd/auth/provider"
)

// user readable errors
var (
	ErrUsernameTaken = errors.New("oopsy, username already taken")
	ErrEmailTaken    = errors.New("oopsy, email is already taken")
	ErrFailedToLogin = errors.New("oopsy, email or password incorrect")
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

func NewAuthService(ap provider.AuthProvider) AuthService {
	return AuthService{authProvider: ap}
}

func (a AuthService) CreateClient(rq auth.RegisterUserRequest) (auth.RegisterUserResponse, error) {
	resp, err := a.authProvider.CreateNewClient(rq.FullName, rq.FullName, rq.Email, rq.Username, rq.Password)
	if err != nil {
		log.Fatalf("failed to create user err: %q", err)
		switch err {
		case provider.ErrSysUsernameTaken:
			return auth.RegisterUserResponse{}, ErrUsernameTaken
		case provider.ErrSysEmailTaken:
			return auth.RegisterUserResponse{}, ErrEmailTaken
		default:
			return auth.RegisterUserResponse{}, nil
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
		log.Fatalf("failed to login err: %q", err)
		switch err {
		case provider.ErrSysFailedToLogin:
			return auth.LoginUserResonse{}, ErrFailedToLogin
		default:
			return auth.LoginUserResonse{}, nil
		}
	}
	return auth.LoginUserResonse{
		JWT:     auth.JWT(resp.JWT),
		Message: SUCCESS_MESSAGE,
	}, nil
}

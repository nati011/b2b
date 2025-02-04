package auth

import (
	"errors"
	"log"

	auth "b2b.nati011.github.com/cmd/auth/model/dto"
	provider "b2b.nati011.github.com/cmd/auth/provider"
)

// user readable errors
var (
	Human_Readable_Error_UsernameTaken = errors.New("Oopsy, username already taken")
	Human_Readable_Error_EmailTaken    = errors.New("Oopsy, email is already taken")
	Human_Readable_Error_FailedToLogin = errors.New("Oopsy, email or password incorrect")
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
		switch err {
		case provider.System_Readable_Error_UsernameTaken:
			return auth.RegisterUserResponse{}, Human_Readable_Error_UsernameTaken
		case provider.System_Readable_Error_EmailTaken:
			return auth.RegisterUserResponse{}, Human_Readable_Error_EmailTaken
		default:
			return auth.RegisterUserResponse{}, nil
		}
		log.Fatalf("failed to create user err: %q", err)
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
		case provider.System_Readable_Error_FailedToLogin:
			return auth.LoginUserResonse{}, Human_Readable_Error_FailedToLogin
		default:
			return auth.LoginUserResonse{}, nil
		}
		log.Fatalf("failed to login err: %q", err)
	}
	return auth.LoginUserResonse{
		JWT:     resp.JWT,
		Message: SUCCESS_MESSAGE,
	}, nil
}

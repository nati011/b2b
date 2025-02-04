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

type Authorizer interface {
	CreateClient(auth.RegisterUserRequest) (auth.RegisterUserResponse, error)
	LoginClient(auth.LoginUserRequest) (auth.LoginUserResonse, error)
}

type AuthService struct {
	authProvider *provider.AuthProvider
}

func NewAuthService(ap *provider.AuthProvider) AuthService {
	return AuthService{authProvider: ap}
}

func (a *AuthService) CreateClient(rq auth.RegisterUserRequest) (auth.RegisterUserResponse, error) {
	resp, err := CreateNewClient(rq.FullName, nil, nil, rq.Email, rq.Username, rq.Password)
	if err != nil {
		switch err {
		case provider.ErrEmailTaken:
			return auth.RegisterUserResponse{}, Human_Readable_Error_UsernameTaken
		case provider.ErrEmailTaken:
			return auth.RegisterUserResponse{}, Human_Readable_Error_EmailTaken
		case provider.ErrFailedToLogin:
		default:
			return auth.RegisterUserResponse{}, nil
		}
		log.Fatalf("failed to create client")
	}
	return resp, nil
}

func (a *AuthService) LoginClient(rq auth.LoginUserRequest) (auth.LoginUserResonse, error) {
	resp, err := a.authProvider
	if err != nil {
		log.Fatalf("failed to create client")
	}
	return resp, nil
}

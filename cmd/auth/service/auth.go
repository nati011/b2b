package auth

import (
	"log"

	auth "b2b.nati011.github.com/cmd/auth/model/dto"
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
	resp, err := a.authProvider.CreateNewClient(rq.FullName, nil, nil, rq.Email, rq.Username, rq.Password)
	if err != nil {
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

package auth

import (
	auth "b2b.nati011.github.com/cmd/auth/model/dto"
	provider "b2b.nati011.github.com/cmd/auth/provider"
)

type Authorizer interface {
	CreateClient(auth.RegisterUserRequest) (auth.RegisterUserResponse, error)
	LoginClient(auth.LoginUserRequest) (auth.LoginUserResonse, error)
}

type AuthService struct {
	authProvider *provider.AuthProvider
}

func NewAuthService() (*AuthService, error) {
	return &AuthService{}, nil
}

func (a *AuthService) CreateClient(auth.RegisterUserRequest) (auth.RegisterUserResponse, error) {
	return auth.RegisterUserResponse{}, nil
}

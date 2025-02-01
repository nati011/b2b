package auth

import (
	auth "b2b.nati011.github.com/cmd/auth/model/dto"
)

type Authorizer interface {
	CreateClient(auth.RegisterUserRequest) (auth.RegisterUserResponse, error)
	LoginClient(auth.LoginUserRequest) (auth.LoginUserResonse, error)
}

type AuthService struct {
}

func NewAuthService() (*AuthService, error) {
	return &AuthService{}, nil
}

func CreateClient(auth.RegisterUserRequest) (auth.RegisterUserResponse, error) {
	return auth.RegisterUserResponse{}, nil
}

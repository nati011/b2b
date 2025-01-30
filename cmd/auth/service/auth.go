package auth

import auth "b2b.nati011.github.com/cmd/auth/model/dto"

type AuthService struct {
}

func (a *AuthService) CreateClient() (auth.RegisterUserResponse, error) {
	return auth.RegisterUserResponse{}, nil
}

package auth

import (
	"b2b.nati011.github.com/internal/core/application/auth"
)

func LoginResponseMapper(in auth.LoginAuthResponse) LoginAuthResponse {
	return (LoginAuthResponse{
		JWT: JWT(in.JWT),
	})
}

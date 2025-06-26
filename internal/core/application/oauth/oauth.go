package oauth

import (
	"context"
	"errors"

	port "b2b.nati011.github.com/internal/port/application/oauth/provider"
)

var (
	ErrFailedToAuthenticate = errors.New("failed to authenticate")
	ErrUnknown              = errors.New("unknown error has occured")
)

type OAuthRequest struct {
	Token string
}

type OAuthResponse struct {
	UserId string
	JWT    JWT
}

type JWT struct {
	AccessToken      string
	IDToken          string
	ExpiresIn        int
	RefreshExpiresIn int
	RefreshToken     string
	TokenType        string
	NotBeforePolicy  int
	SessionState     string
	Scope            string
}

type Provider interface {
	GoogleSignOn(ctx context.Context, Request OAuthRequest) (OAuthResponse, error)
}

type OAuthService struct {
	oauthProvider port.Provider
}

func NewAuthService(
	AP port.Provider,
) Provider {
	return &OAuthService{
		oauthProvider: AP,
	}
}

func (o *OAuthService) GoogleSignOn(ctx context.Context, Request OAuthRequest) (OAuthResponse, error) {
	resp, err := o.oauthProvider.GoogleSignOn(ctx, &port.OAuthRequest{
		Token: Request.Token,
	})
	if err != nil {
		switch err {
		case port.ErrSysFailedToLogin:
			return OAuthResponse{}, ErrFailedToAuthenticate
		default:
			return OAuthResponse{}, ErrUnknown
		}
	}

	return OAuthResponse{
		JWT: JWT(resp.JWT),
	}, nil
}

package oauth

import (
	"context"
	"errors"

	"b2b.nati011.github.com/internal/core/application/event"
	"b2b.nati011.github.com/internal/core/util"
	auth_port "b2b.nati011.github.com/internal/port/application/auth/provider"
	port "b2b.nati011.github.com/internal/port/application/oauth/provider"
)

var (
	ErrFailedToAuthenticate = errors.New("failed to authenticate")
	ErrUnknown              = errors.New("unknown error has occured")
)

type OAuthRequest struct {
	Token       string `json:"token"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}

type OAuthResponse struct {
	JWT JWT
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
	authProvider  auth_port.Provider
	event         *event.Broker
}

func NewAuthService(
	op port.Provider,
	e *event.Broker,
	ap auth_port.Provider,

) Provider {
	return &OAuthService{
		oauthProvider: op,
		event:         e,
		authProvider:  ap,
	}
}

func (o *OAuthService) GoogleSignOn(ctx context.Context, request OAuthRequest) (OAuthResponse, error) {
	resp, err := o.oauthProvider.GoogleSignOn(ctx, &port.OAuthRequest{
		Token: request.Token,
	})
	if err != nil {
		switch err {
		case port.ErrSysFailedToLogin:
			return OAuthResponse{}, ErrFailedToAuthenticate
		default:
			return OAuthResponse{}, ErrUnknown
		}
	}

	payload := util.EventRetailerSSOPayload{
		ID:        resp.UserId,
		Email:     request.Email,
		Username:  request.Email,
		FirstName: request.FirstName,
		LastName:  request.LastName,
	}

	o.event.Publish(util.EVENT_RETAILER_SSO, payload)

	return OAuthResponse{
		JWT: JWT(resp.JWT),
	}, nil
}

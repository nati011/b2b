package provider

import (
	"context"
	"errors"
)

// system errors
var (
	ErrSysUnknown        = errors.New("unknown error")
	ErrSysFailedToLogin  = errors.New("authentication failed")
	ErrUserAlreadyExists = errors.New("user already exists")
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
	GoogleSignOn(ctx context.Context, Request *OAuthRequest) (OAuthResponse, error)
}

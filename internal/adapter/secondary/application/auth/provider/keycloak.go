package provider

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/Nerzal/gocloak/v13"

	port "b2b.nati011.github.com/internal/port/application/auth/provider"
)

const (
	MessageErrKeyCloakEmailTaken    = "User exists with same email"
	MessageErrKeyCloakUsernameTaken = "User exists with same username"

	MessageErrFailedLogin = "Invalid user credentials"

	// KeycloakInstanceURL      = "http://localhost:8080"
	// KeycloakUsername         = "admin"
	// KeycloakPassword         = "admin"
	// KeycloakRealm            = "master"
	// KeycloakApplicationRealm = "test"
	// KeycloakClientId         = "test"
	// KeycloakClientSecret     = "jQILbkSn6ywmVVYjxgSMHYOUfnlA7pMS"
)

type KeycloakProvider struct {
	KeycloakInstanceURL      string
	KeycloakUsername         string
	KeycloakPassword         string
	KeycloakRealm            string
	KeycloakApplicationRealm string
	KeycloakClientId         string
	KeycloakClientSecret     string
}

func NewKeycloakProvider(
	keycloakInstanceURL string,
	keycloakUsername string,
	keycloakPassword string,
	keycloakRealm string,
	keycloakApplicationRealm string,
	keycloakClientId string,
	keycloakClientSecret string,
) port.Provider {
	return &KeycloakProvider{
		KeycloakInstanceURL:      keycloakInstanceURL,
		KeycloakUsername:         keycloakUsername,
		KeycloakPassword:         keycloakPassword,
		KeycloakRealm:            keycloakRealm,
		KeycloakApplicationRealm: keycloakApplicationRealm,
		KeycloakClientId:         keycloakClientId,
		KeycloakClientSecret:     keycloakClientSecret,
	}
}

func (k KeycloakProvider) CreateNewClient(ctx context.Context, req port.RegisterUserRequest) (port.RegisterUserResponse, error) {
	client := gocloak.NewClient(k.KeycloakInstanceURL)

	token, err := client.LoginAdmin(ctx, k.KeycloakUsername, k.KeycloakPassword, k.KeycloakRealm)
	if err != nil {
		log.Printf("Something wrong with the credentials or URL: %v", err)
		return port.RegisterUserResponse{}, port.ErrSysUnknown
	}

	user := gocloak.User{
		FirstName: gocloak.StringP(req.FirstName),
		LastName:  gocloak.StringP(req.LastName),
		Email:     gocloak.StringP(req.Email),
		Enabled:   gocloak.BoolP(true),
		Username:  gocloak.StringP(req.Username),
		ID:        gocloak.StringP(req.Email),
	}

	userId, err := client.CreateUser(ctx, token.AccessToken, k.KeycloakApplicationRealm, user)
	if err != nil {
		var apiErr *gocloak.APIError
		if errors.As(err, &apiErr) {
			switch apiErr.Code {
			case 409:
				switch {
				case strings.Contains(apiErr.Message, MessageErrKeyCloakEmailTaken):
					return port.RegisterUserResponse{}, port.ErrSysEmailTaken
				case strings.Contains(apiErr.Message, MessageErrKeyCloakUsernameTaken):
					return port.RegisterUserResponse{}, port.ErrSysUsernameTaken
				}
			default:
				return port.RegisterUserResponse{}, err

			}
		}
	}

	if req.Email != "" {
		client.SendVerifyEmail(ctx, token.AccessToken, userId, k.KeycloakApplicationRealm)
	}
	if req.Password != "" {
		client.SetPassword(ctx, token.AccessToken, userId, k.KeycloakApplicationRealm, req.Password, false)
	}

	if err != nil {
		return port.RegisterUserResponse{}, err
	}

	return port.RegisterUserResponse{
		Id:       userId,
		Username: req.Username,
	}, nil
}

func (k KeycloakProvider) DeleteClient(ctx context.Context, userId string) error {
	client := gocloak.NewClient(k.KeycloakInstanceURL)
	token, err := client.LoginAdmin(ctx, k.KeycloakUsername, k.KeycloakPassword, k.KeycloakRealm)
	if err != nil {
		log.Printf("Something wrong with the credentials or URL: %v", err)
		return port.ErrSysUnknown
	}
	err = client.DeleteUser(ctx, token.AccessToken, k.KeycloakRealm, userId)
	if err != nil {
		return port.ErrSysUnknown
	}
	return nil
}

func (k KeycloakProvider) ClientLogin(ctx context.Context, req port.LoginUserRequest) (port.LoginAuthResponse, error) {
	client := gocloak.NewClient(k.KeycloakInstanceURL)

	token, err := client.Login(ctx, k.KeycloakClientId, k.KeycloakClientSecret, k.KeycloakApplicationRealm, req.Email, req.Password)

	if err != nil {
		var apiErr *gocloak.APIError
		if errors.As(err, &apiErr) {
			switch apiErr.Code {
			case 401:
				switch {
				case strings.Contains(apiErr.Message, MessageErrFailedLogin):
					return port.LoginAuthResponse{}, port.ErrSysFailedToLogin
				}
			default:
				return port.LoginAuthResponse{}, port.ErrSysUnknown

			}
		}
	}

	if token == nil {
		return port.LoginAuthResponse{}, port.ErrSysFailedToLogin
	}

	return port.LoginAuthResponse{
		JWT: port.JWT{
			AccessToken:      token.AccessToken,
			IDToken:          token.IDToken,
			ExpiresIn:        token.ExpiresIn,
			RefreshExpiresIn: token.RefreshExpiresIn,
			RefreshToken:     token.RefreshToken,
			TokenType:        token.TokenType,
			NotBeforePolicy:  token.NotBeforePolicy,
			SessionState:     token.SessionState,
			Scope:            token.Scope,
		},
	}, err
}

func (k *KeycloakProvider) RefreshToken(ctx context.Context, req port.RefreshTokenRequest) (port.LoginAuthResponse, error) {
	client := gocloak.NewClient(k.KeycloakInstanceURL)

	token, err := client.RefreshToken(ctx, req.RefreshToken, k.KeycloakClientId, k.KeycloakClientSecret, k.KeycloakRealm)
	if err != nil {
		var apiErr *gocloak.APIError
		if errors.As(err, &apiErr) {
			switch apiErr.Code {
			case 401:
				switch {
				case strings.Contains(apiErr.Message, MessageErrFailedLogin):
					return port.LoginAuthResponse{}, port.ErrSysFailedToLogin
				}
			default:
				return port.LoginAuthResponse{}, port.ErrSysFailedToLogin

			}
		}
	}

	return port.LoginAuthResponse{
		JWT: port.JWT{
			AccessToken:      token.AccessToken,
			IDToken:          token.IDToken,
			ExpiresIn:        token.ExpiresIn,
			RefreshExpiresIn: token.RefreshExpiresIn,
			RefreshToken:     token.RefreshToken,
			TokenType:        token.TokenType,
			NotBeforePolicy:  token.NotBeforePolicy,
			SessionState:     token.SessionState,
			Scope:            token.Scope,
		},
	}, err
}

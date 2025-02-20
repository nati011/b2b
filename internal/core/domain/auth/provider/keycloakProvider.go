package provider

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/Nerzal/gocloak/v13"
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
) *KeycloakProvider {
	return &KeycloakProvider{
		KeycloakInstanceURL:      keycloakInstanceURL,
		KeycloakUsername:         keycloakUsername,
		KeycloakPassword:         keycloakPassword,
		KeycloakRealm:            keycloakRealm,
		KeycloakApplicationRealm: keycloakApplicationRealm,
		KeycloakClientId:         keycloakClientId,
	}
}

func (k KeycloakProvider) CreateNewClient(firstName string, lastName string, email string, username string, password string) (CreateClientAuthResponse, error) {
	client := gocloak.NewClient(k.KeycloakInstanceURL)
	ctx := context.Background()

	token, err := client.LoginAdmin(ctx, k.KeycloakUsername, k.KeycloakPassword, k.KeycloakRealm)
	if err != nil {
		log.Fatalf("Something wrong with the credentials or URL: %v", err)
	}

	user := gocloak.User{
		FirstName: gocloak.StringP(firstName),
		LastName:  gocloak.StringP(lastName),
		Email:     gocloak.StringP(email),
		Enabled:   gocloak.BoolP(true),
		Username:  gocloak.StringP(username),
		ID:        gocloak.StringP(email),
	}

	userId, err := client.CreateUser(ctx, token.AccessToken, k.KeycloakApplicationRealm, user)
	if err != nil {
		var apiErr *gocloak.APIError
		if errors.As(err, &apiErr) {
			switch apiErr.Code {
			case 409:
				switch {
				case strings.Contains(apiErr.Message, MessageErrKeyCloakEmailTaken):
					return CreateClientAuthResponse{}, ErrSysEmailTaken
				case strings.Contains(apiErr.Message, MessageErrKeyCloakUsernameTaken):
					return CreateClientAuthResponse{}, ErrSysUsernameTaken
				}
			default:
				return CreateClientAuthResponse{}, ErrSysUnknown

			}
		}
	}

	if email != "" {
		client.SendVerifyEmail(ctx, token.AccessToken, userId, k.KeycloakApplicationRealm)
	}
	if password != "" {
		client.SetPassword(ctx, token.AccessToken, userId, k.KeycloakApplicationRealm, password, false)
	}

	if err != nil {
		return CreateClientAuthResponse{}, ErrSysUnknown
	}

	return CreateClientAuthResponse{
		Username: username,
	}, nil
}

func (k KeycloakProvider) ClientLogin(email, password string) (LoginAuthResponse, error) {
	client := gocloak.NewClient(k.KeycloakInstanceURL)
	ctx := context.Background()
	adminToken, err := client.LoginAdmin(ctx, k.KeycloakUsername, k.KeycloakPassword, k.KeycloakRealm)
	if err != nil {
		log.Fatalf("Something wrong with the credentials or URL: %v", err)
	}
	clientSecret, err := client.GetClientSecret(ctx, adminToken.AccessToken, k.KeycloakApplicationRealm, k.KeycloakClientId)
	if err != nil {
		log.Fatal("client secret fetching failed:" + err.Error())
	}
	token, err := client.Login(ctx, k.KeycloakClientId, *clientSecret.Value, k.KeycloakApplicationRealm, email, password)
	if err != nil {
		var apiErr *gocloak.APIError
		if errors.As(err, &apiErr) {
			switch apiErr.Code {
			case 401:
				return LoginAuthResponse{}, ErrSysFailedToLogin

			default:
				return LoginAuthResponse{}, ErrSysFailedToLogin

			}
		}
	}
	return LoginAuthResponse{
		JWT: JWT{
			token.AccessToken,
			token.IDToken,
			token.ExpiresIn,
			token.RefreshExpiresIn,
			token.RefreshToken,
			token.TokenType,
			token.NotBeforePolicy,
			token.SessionState,
			token.Scope,
		},
	}, err
}

func (k KeycloakProvider) RefreshToken(refreshToken string) (LoginAuthResponse, error) {
	client := gocloak.NewClient(k.KeycloakInstanceURL)
	ctx := context.Background()
	adminToken, err := client.LoginAdmin(ctx, k.KeycloakUsername, k.KeycloakPassword, k.KeycloakRealm)
	if err != nil {
		log.Fatalf("Something wrong with the credentials or URL: %v", err)
	}
	clientSecret, err := client.GetClientSecret(ctx, adminToken.AccessToken, k.KeycloakApplicationRealm, k.KeycloakClientId)
	if err != nil {
		log.Fatal("Client secret fetching failed:" + err.Error())
	}
	token, err := client.RefreshToken(ctx, k.KeycloakClientId, *clientSecret.Value, k.KeycloakApplicationRealm, refreshToken)
	if err != nil {
		var apiErr *gocloak.APIError
		if errors.As(err, &apiErr) {
			switch apiErr.Code {
			case 401:
				switch {
				case strings.Contains(apiErr.Message, MessageErrFailedLogin):
					return LoginAuthResponse{}, ErrSysFailedToLogin
				}
			default:
				return LoginAuthResponse{}, ErrSysFailedToLogin

			}
		}
	}

	return LoginAuthResponse{
		JWT: JWT{
			token.AccessToken,
			token.IDToken,
			token.ExpiresIn,
			token.RefreshExpiresIn,
			token.RefreshToken,
			token.TokenType,
			token.NotBeforePolicy,
			token.SessionState,
			token.Scope,
		},
	}, err
}

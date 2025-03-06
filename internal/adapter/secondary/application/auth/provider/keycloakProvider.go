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
) port.AuthProvider {
	return &KeycloakProvider{
		KeycloakInstanceURL:      keycloakInstanceURL,
		KeycloakUsername:         keycloakUsername,
		KeycloakPassword:         keycloakPassword,
		KeycloakRealm:            keycloakRealm,
		KeycloakApplicationRealm: keycloakApplicationRealm,
		KeycloakClientId:         keycloakClientId,
	}
}

func (k KeycloakProvider) CreateNewClient(firstName string, lastName string, email string, username string, password string) (port.CreateClientAuthResponse, error) {
	client := gocloak.NewClient(k.KeycloakInstanceURL)
	ctx := context.Background()

	token, err := client.LoginAdmin(ctx, k.KeycloakUsername, k.KeycloakPassword, k.KeycloakRealm)
	if err != nil {
		log.Printf("Something wrong with the credentials or URL: %v", err)
		return port.CreateClientAuthResponse{}, port.ErrSysUnknown
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
			print(apiErr.Message)
			switch apiErr.Code {
			case 409:
				switch {
				case strings.Contains(apiErr.Message, MessageErrKeyCloakEmailTaken):
					return port.CreateClientAuthResponse{}, port.ErrSysEmailTaken
				case strings.Contains(apiErr.Message, MessageErrKeyCloakUsernameTaken):
					return port.CreateClientAuthResponse{}, port.ErrSysUsernameTaken
				}
			default:
				return port.CreateClientAuthResponse{}, port.ErrSysUnknown

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
		return port.CreateClientAuthResponse{}, port.ErrSysUnknown
	}

	return port.CreateClientAuthResponse{
		Username: username,
	}, nil
}

func (k KeycloakProvider) ClientLogin(email, password string) (port.LoginAuthResponse, error) {
	client := gocloak.NewClient(k.KeycloakInstanceURL)
	ctx := context.Background()
	adminToken, err := client.LoginAdmin(ctx, k.KeycloakUsername, k.KeycloakPassword, k.KeycloakRealm)
	if err != nil {
		log.Printf("Something wrong with the credentials or URL: %v", err)
		return port.LoginAuthResponse{}, port.ErrSysUnknown
	}
	clientSecret, err := client.GetClientSecret(ctx, adminToken.AccessToken, k.KeycloakApplicationRealm, k.KeycloakClientId)
	if err != nil {
		log.Printf("client secret fetching failed: %v", err.Error())
		return port.LoginAuthResponse{}, port.ErrSysUnknown
	}
	token, err := client.Login(ctx, k.KeycloakClientId, *clientSecret.Value, k.KeycloakApplicationRealm, email, password)
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

	rptResult, err := client.RetrospectToken(ctx, token.AccessToken, k.KeycloakClientId, k.KeycloakClientSecret, k.KeycloakApplicationRealm)
	if err != nil {
		log.Fatal("Inspection failed:" + err.Error())
		return port.LoginAuthResponse{}, err
	}

	if !*rptResult.Active {
		err := errors.New("token is not active")
		log.Fatal("token is not active:" + err.Error())
		return port.LoginAuthResponse{}, err
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
func (k *KeycloakProvider) RefreshToken(refreshToken string) (port.LoginAuthResponse, error) {
	client := gocloak.NewClient(k.KeycloakInstanceURL)
	ctx := context.Background()
	adminToken, err := client.LoginAdmin(ctx, k.KeycloakUsername, k.KeycloakPassword, k.KeycloakRealm)
	if err != nil {
		log.Printf("Something wrong with the credentials or URL: %v", err)
		return port.LoginAuthResponse{}, port.ErrSysUnknown
	}
	clientSecret, err := client.GetClientSecret(ctx, adminToken.AccessToken, k.KeycloakApplicationRealm, k.KeycloakClientId)
	if err != nil {
		log.Printf("client secret fetching failed: %v", err.Error())
		return port.LoginAuthResponse{}, port.ErrSysUnknown
	}
	token, err := client.RefreshToken(ctx, refreshToken, k.KeycloakClientId, *clientSecret.Value, k.KeycloakRealm)
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

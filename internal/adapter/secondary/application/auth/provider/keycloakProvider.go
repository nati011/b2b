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

func (k KeycloakProvider) CreateNewClient(firstName string, lastName string, email string, username string, password string) (port.CreateClientAuthResonse, error) {
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
					return port.CreateClientAuthResonse{}, port.ErrSysEmailTaken
				case strings.Contains(apiErr.Message, MessageErrKeyCloakUsernameTaken):
					return port.CreateClientAuthResonse{}, port.ErrSysUsernameTaken
				}
			default:
				return port.CreateClientAuthResonse{}, port.ErrSysUnknown

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
		return port.CreateClientAuthResonse{}, port.ErrSysUnknown
	}

	return port.CreateClientAuthResonse{
		Username: username,
	}, nil
}

func (k KeycloakProvider) ClientLogin(email, password string) (port.LoginAuthResonse, error) {
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
				switch {
				case strings.Contains(apiErr.Message, MessageErrFailedLogin):
					return port.LoginAuthResonse{}, port.ErrSysFailedToLogin
				}
			default:
				return port.LoginAuthResonse{}, port.ErrSysFailedToLogin

			}
		}
	}

	rptResult, err := client.RetrospectToken(ctx, token.AccessToken, k.KeycloakClientId, k.KeycloakClientSecret, k.KeycloakApplicationRealm)
	if err != nil {
		log.Fatal("Inspection failed:" + err.Error())
		return port.LoginAuthResonse{}, err
	}

	if !*rptResult.Active {
		err := errors.New("token is not active")
		log.Fatal("token is not active:" + err.Error())
		return port.LoginAuthResonse{}, err
	}

	return port.LoginAuthResonse{
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

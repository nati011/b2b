package provider

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/Nerzal/gocloak/v13"
)

const (
	MessageErrKeyCloakUsernameTaken = "User exists with same email"
	MessageErrKeyCloakEmailTaken    = "User exists with same username"
	MessageErrFailedLogin           = "Invalid user credentials"

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
) *KeycloakProvider {
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

func (k *KeycloakProvider) CreateNewClient(firstName string, lastName string, email string, username string, password string) (string, error) {
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
					return "", ErrSysEmailTaken
				case strings.Contains(apiErr.Message, MessageErrKeyCloakUsernameTaken):
					return "", ErrSysUsernameTaken
				}
			default:
				return "", ErrSysUnknown

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
		return "", ErrSysUnknown
	}

	return username, nil
}

func (k *KeycloakProvider) ClientLogin(email, password string) (string, error) {
	client := gocloak.NewClient(k.KeycloakInstanceURL)
	ctx := context.Background()

	token, err := client.Login(ctx, k.KeycloakClientId, k.KeycloakClientSecret, k.KeycloakApplicationRealm, email, password)
	if err != nil {
		var apiErr *gocloak.APIError
		if errors.As(err, &apiErr) {
			switch apiErr.Code {
			case 401:
				switch {
				case strings.Contains(apiErr.Message, MessageErrFailedLogin):
					return "", ErrSysFailedToLogin
				}
			default:
				return "", ErrSysFailedToLogin

			}
		}
	}

	rptResult, err := client.RetrospectToken(ctx, token.AccessToken, k.KeycloakClientId, k.KeycloakClientSecret, k.KeycloakApplicationRealm)
	if err != nil {
		log.Fatal("Inspection failed:" + err.Error())
		return "", err
	}

	if !*rptResult.Active {
		err := errors.New("token is not active")
		log.Fatal("token is not active:" + err.Error())
		return "", err
	}

	return token.AccessToken, err
}

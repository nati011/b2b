package provider

import (
	"context"
	"errors"
	"log"

	"github.com/Nerzal/gocloak/v13"
)

const (
	KeycloakInstanceURL      = "http://localhost:8080"
	KeycloakUsername         = "admin"
	KeycloakPassword         = "admin"
	KeycloakRealm            = "master"
	KeycloakApplicationRealm = "test"
	KeycloakClientId         = "test"
	KeycloakClientSecret     = "jQILbkSn6ywmVVYjxgSMHYOUfnlA7pMS"
)

type (
	KeycloakProvider struct{}
)

func NewKeycloakProvider() *KeycloakProvider {
	return &KeycloakProvider{}
}

func (k *KeycloakProvider) CreateNewClient(firstName string, lastName string, email string, username string, password string) (string, error) {
	client := gocloak.NewClient(KeycloakInstanceURL)
	ctx := context.Background()

	token, err := client.LoginAdmin(ctx, KeycloakUsername, KeycloakPassword, KeycloakRealm)
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

	userId, err := client.CreateUser(ctx, token.AccessToken, KeycloakApplicationRealm, user)
	client.SendVerifyEmail(ctx, token.AccessToken, userId, KeycloakApplicationRealm)
	client.SetPassword(ctx, token.AccessToken, userId, KeycloakApplicationRealm, password, false)

	if err != nil {
		log.Fatalf("Oh no!, failed to create user: %v", err)
	}

	return userId, nil
}

func (k *KeycloakProvider) ClientLogin(email, password string) (string, error) {
	client := gocloak.NewClient(KeycloakInstanceURL)
	ctx := context.Background()
	token, err := client.Login(ctx, KeycloakClientId, KeycloakClientSecret, KeycloakApplicationRealm, email, password)
	if err != nil {
		return "", err
	}

	rptResult, err := client.RetrospectToken(ctx, token.AccessToken, KeycloakClientId, KeycloakClientSecret, KeycloakApplicationRealm)
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

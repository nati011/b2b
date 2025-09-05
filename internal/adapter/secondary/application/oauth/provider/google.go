package oauth

import (
	"context"

	port "b2b.nati011.github.com/internal/port/application/oauth/provider"
	"b2b.nati011.github.com/package/keycloak"
)

type KeycloakOAuthProvider struct {
	KeycloakInstanceURL  string
	KeycloakRealm        string
	KeycloakClientId     string
	KeycloakClientSecret string
}

func NewKeycloakOAuthProvider(
	keycloakInstanceURL string,
	keycloakRealm string,
	keycloakClientId string,
	keycloakClientSecret string,
) port.Provider {
	return &KeycloakOAuthProvider{
		KeycloakInstanceURL:  keycloakInstanceURL,
		KeycloakRealm:        keycloakRealm,
		KeycloakClientId:     keycloakClientId,
		KeycloakClientSecret: keycloakClientSecret,
	}
}

func (k *KeycloakOAuthProvider) GoogleSignOn(ctx context.Context, request *port.OAuthRequest) (port.OAuthResponse, error) {
	client := keycloak.NewClient(k.KeycloakInstanceURL, k.KeycloakRealm)
	token, err := client.GetToken(ctx, k.KeycloakClientId, k.KeycloakClientSecret, request.Token, "google")

	if err != nil {
		switch err {
		case keycloak.ErrUserAlreadyExists:
			return port.OAuthResponse{}, port.ErrUserAlreadyExists
		default:
			return port.OAuthResponse{}, port.ErrSysFailedToLogin
		}
	}

	return port.OAuthResponse{
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

package keycloak

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
)

type KeyCloakAdapter struct {
	basePath      string
	realm         string
	restyClient   *resty.Client
	tokenEndpoint string
}

type JWT struct {
	AccessToken      string `json:"access_token"`
	IDToken          string `json:"id_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	NotBeforePolicy  int    `json:"not-before-policy"`
	SessionState     string `json:"session_state"`
	Scope            string `json:"scope"`
}

type TokenOptions struct {
	GrantType        string `json:"grant_type,omitempty"`
	Scope            string `json:"scope,omitempty"`
	SubjectToken     string `json:"subject_token,omitempty"`
	SubjectTokenType string `json:"subject_token_type,omitempty"`
	Audience         string `json:"audience,omitempty"`
	SubjectIssuer    string `json:"subject_issuer"`
}

func NewClient(basePath string, realm string) *KeyCloakAdapter {
	c := KeyCloakAdapter{
		basePath:    strings.TrimRight(basePath, "/"),
		restyClient: resty.New(),
		realm:       realm,
	}
	path := []string{
		c.basePath,
		"realms",
		c.realm,
		"protocol",
		"openid-connect",
		"token",
	}
	c.tokenEndpoint = strings.Join(path, "/")
	return &c
}

func (k *KeyCloakAdapter) GetRequestWithBasicAuth(ctx context.Context, clientID, clientSecret string) *resty.Request {
	req := k.restyClient.R().SetContext(ctx)
	req.SetHeader("Content-Type", "application/x-www-form-urlencoded")
	httpBasicAuth := base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))
	req.SetHeader("Authorization", "Basic "+httpBasicAuth)
	return req
}

func (k *KeyCloakAdapter) GetToken(ctx context.Context, clientID, clientSecret, subjectToken, subjectIssuer string) (*JWT, error) {
	tokenOptions := TokenOptions{
		GrantType:        "urn:ietf:params:oauth:grant-type:token-exchange",
		SubjectToken:     subjectToken,
		SubjectTokenType: "urn:ietf:params:oauth:token-type:access_token",
		SubjectIssuer:    subjectIssuer,
		Scope:            "openid",
	}

	req := k.GetRequestWithBasicAuth(ctx, clientID, clientSecret)

	var token JWT
	resp, err := req.
		SetFormData(map[string]string{
			"grant_type":         tokenOptions.GrantType,
			"subject_token":      tokenOptions.SubjectToken,
			"subject_token_type": tokenOptions.SubjectTokenType,
			"subject_issuer":     tokenOptions.SubjectIssuer,
			"scope":              tokenOptions.Scope,
		}).
		SetResult(&token).
		Post(k.tokenEndpoint)

	if err != nil {
		return nil, fmt.Errorf("failed to make token request: %w", err)
	}

	if resp.IsError() {
		if string(resp.Body()) == "{\"error\":\"invalid_token\",\"error_description\":\"User already exists\"}" {
			return nil, fmt.Errorf("token request failed with status %d: email: %v already in use", resp.StatusCode(), resp.String())
		}
		return nil, fmt.Errorf("token request failed with status %d: %s", resp.StatusCode(), resp.String())
	}

	return &token, nil
}

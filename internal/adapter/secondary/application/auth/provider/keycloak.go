package provider

import (
	"context"
	"errors"
	"log"
	"strings"

	port "b2b.nati011.github.com/internal/port/application/auth/provider"
	"github.com/Nerzal/gocloak/v13"
	"github.com/golang-jwt/jwt"
)

const (
	MessageErrKeyCloakEmailTaken    = "User exists with same email"
	MessageErrKeyCloakUsernameTaken = "User exists with same username"
	MessageErrFailedLogin           = "Invalid user credentials"
)

var (
	ErrUnknown = errors.New("¯\\_(o_o)_/¯, unknown error")
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

func (k *KeycloakProvider) RetrospectToken(ctx context.Context, token string) (port.RetrospectionResult, error) {
	client := gocloak.NewClient(k.KeycloakInstanceURL)
	result, err := client.RetrospectToken(
		ctx,
		token,
		k.KeycloakClientId,
		k.KeycloakClientSecret,
		k.KeycloakRealm,
	)

	return port.RetrospectionResult{
		Exp:      result.Exp,
		Nbf:      result.Nbf,
		Iat:      result.Iat,
		Active:   result.Active,
		AuthTime: result.AuthTime,
		Jti:      result.Jti,
		Type:     result.Type,
	}, err
}

func GetEmail(c jwt.MapClaims) (string, error) {
	var cs []string
	switch v := c["email"].(type) {
	case string:
		cs = append(cs, v)
	case []string:
		cs = v
	case []interface{}:
		for _, a := range v {
			vs, ok := a.(string)
			if !ok {
				return "", ErrUnknown
			}
			cs = append(cs, vs)
		}
	}
	return strings.Join(cs, ""), nil
}

func (k *KeycloakProvider) DecodeToken(ctx context.Context, token string) (port.DecodedResult, error) {
	client := gocloak.NewClient(k.KeycloakInstanceURL)
	decodedToken, mapClaims, err := client.DecodeAccessToken(
		ctx,
		token,
		k.KeycloakApplicationRealm,
	)
	if err != nil {
		log.Printf("failed to decode token: %v", err)
		return port.DecodedResult{}, port.ErrSysUnknown
	}
	res := port.DecodedResult{
		Raw:       decodedToken.Raw,
		Header:    decodedToken.Header,
		Signature: decodedToken.Signature,
		Valid:     decodedToken.Valid}
	claims := decodedToken.Claims
	expirationTime, err := claims.GetExpirationTime()
	if err != nil {
		log.Printf("failed to decode token: %v", err)
		return port.DecodedResult{}, port.ErrSysUnknown
	}
	issueDate, err := claims.GetIssuedAt()
	if err != nil {
		log.Printf("failed to decode token: %v", err)
		return port.DecodedResult{}, port.ErrSysUnknown
	}
	subject, err := claims.GetSubject()
	if err != nil {
		log.Printf("failed to decode token: %v", err)
		return port.DecodedResult{}, port.ErrSysUnknown
	}
	issuer, err := claims.GetIssuer()
	if err != nil {
		log.Printf("failed to decode token: %v", err)
		return port.DecodedResult{}, port.ErrSysUnknown
	}
	email, err := GetEmail(jwt.MapClaims(*mapClaims))
	if err != nil {
		log.Printf("failed to decode token: %v", err)
		return port.DecodedResult{}, port.ErrSysUnknown
	}
	res.Claims = port.Claims{
		ExpirationTime: expirationTime.Time,
		IssuedAt:       issueDate.Time,
		Issuer:         issuer,
		Subject:        subject,
		Email:          email,
	}
	return res, nil
}
func (k *KeycloakProvider) CreateNewClient(ctx context.Context, req *port.RegisterUserRequest) (port.RegisterUserResponse, error) {
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

func (k *KeycloakProvider) DeleteClient(ctx context.Context, userId string) error {
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

func (k *KeycloakProvider) ResetPassword(ctx context.Context, userId, new_password string) error {
	client := gocloak.NewClient(k.KeycloakInstanceURL)
	token, err := client.LoginAdmin(ctx, k.KeycloakUsername, k.KeycloakPassword, k.KeycloakRealm)
	if err != nil {
		log.Printf("Something wrong with the credentials or URL: %v", err)
		return port.ErrSysUnknown
	}
	err = client.SetPassword(ctx, token.AccessToken, userId, k.KeycloakRealm, new_password, false)
	if err != nil {
		return port.ErrSysUnknown
	}
	return nil
}

func (k KeycloakProvider) ClientLogin(ctx context.Context, req *port.LoginUserRequest) (port.LoginAuthResponse, error) {
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

func (k KeycloakProvider) ClientLogout(ctx context.Context, req port.RefreshTokenRequest) error {
	client := gocloak.NewClient(k.KeycloakInstanceURL)
	err := client.Logout(ctx, k.KeycloakClientId, k.KeycloakClientSecret, k.KeycloakRealm, req.RefreshToken)
	if err != nil {
		return port.ErrSysUnknown
	}
	return nil
}

func (k *KeycloakProvider) RefreshToken(ctx context.Context, req *port.RefreshTokenRequest) (port.LoginAuthResponse, error) {
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

func (k *KeycloakProvider) AssignRole(ctx context.Context, UserId string, roleId int) error {
	// client := gocloak.NewClient(k.KeycloakInstanceURL)
	// client.AddClientRolesToUser(ctx, client.RestyClient().Token, )
	return nil
}

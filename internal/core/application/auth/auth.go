package auth

import (
	"context"
	"crypto/rand"
	"errors"
	"log"
	"math/big"
	"time"

	"b2b.nati011.github.com/internal/core/application/email"
	"b2b.nati011.github.com/internal/core/application/role"
	port "b2b.nati011.github.com/internal/port/application/auth/provider"
	"github.com/golang-jwt/jwt"
)

var (
	ErrUsernameTaken        = errors.New("¯\\_(o_o)_/¯, username already taken")
	ErrEmailTaken           = errors.New("¯\\_(o_o)_/¯, email already taken")
	ErrFailedToLogin        = errors.New("¯\\_(o_o)_/¯, email or password incorrect")
	ErrUnknown              = errors.New("¯\\_(o_o)_/¯, unknown error has occured")
	ErrEmailNotSupplied     = errors.New("¯\\_(o_o)_/¯, email mandatory")
	ErrUsernameNotSupplied  = errors.New("¯\\_(o_o)_/¯, username mandatory")
	ErrInvalidEmail         = errors.New("¯\\_(o_o)_/¯, invalid Email")
	ErrPasswordNotSupplied  = errors.New("¯\\_(o_o)_/¯, password mandatory")
	ErrFirstNameNotSupplied = errors.New("¯\\_(o_o)_/¯, firstname mandatory")
	ErrLastNameNotSupplied  = errors.New("¯\\_(o_o)_/¯, lastname mandatory")
	ErrTokenNotSupplied     = errors.New("¯\\_(o_o)_/¯, token is mandatory")
	ErrTokenHasExpired      = errors.New("¯\\_(o_o)_/¯, token has expired")
	ErrUserIdNotSupplied    = errors.New("¯\\_(o_o)_/¯, userId mandatory")
	ErrFailedToDecodeToken  = errors.New("¯\\_(o_o)_/¯, failed to decode token")
)

type RegisterUserRequest struct {
	Email       string
	Password    string
	BirthDate   time.Time
	PhoneNumber string
	ExternalId  string
	FirstName   string
	LastName    string
	Username    string
}

type RegisterUserWithoutPasswordRequest struct {
	Email       string
	BirthDate   time.Time
	PhoneNumber string
	ExternalId  string
	FirstName   string
	LastName    string
	Username    string
}

type RegisterUserResponse struct {
	Id       string
	Username string
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type LoginAuthResponse struct {
	JWT JWT
}

type LoginUserRequest struct {
	Email    string
	Password string
}

type ResetCredentialsRequest struct {
	NewPassword string
	ResetToken  string
}

type InitClientCredentialsResetRequest struct {
	UserId string `json:"user_id"`
	Email  string `json:"email"`
}

type JWT struct {
	AccessToken      string `json:"access_token"`
	IDToken          string `json:"id_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	NotBeforePolicy  int    `json:"not_before_policy"`
	SessionState     string `json:"session_state"`
	Scope            string `json:"scope"`
}

type RetrospectionResult struct {
	Active bool `json:"active"`
}

type Claims struct {
	ExpirationTime time.Time
	IssuedAt       time.Time
	NotBefore      time.Time
	Issuer         string
	Subject        string
	Audience       string
	Email          string
}

type DecodeResult struct {
	Claims Claims `json:"claims"`
}

type ResourceAccess struct {
	Roles []string `json:"roles"`
}

type Provider interface {
	CreateNewClientWithPassword(ctx context.Context, req *RegisterUserRequest) (RegisterUserResponse, error)
	CreateNewClientWithOutPassword(ctx context.Context, req *RegisterUserWithoutPasswordRequest) (RegisterUserResponse, error)
	ClientLogin(ctx context.Context, req *LoginUserRequest) (LoginAuthResponse, error)
	RefreshToken(ctx context.Context, req *RefreshTokenRequest) (LoginAuthResponse, error)
	DeleteClient(ctx context.Context, userId string) error
	ResetClientCredentials(ctx context.Context, req *ResetCredentialsRequest) error
	RetrospectToken(ctx context.Context, token string) (RetrospectionResult, error)
	DecodeToken(ctx context.Context, token string) (DecodeResult, error)
	ClientLogout(ctx context.Context, req RefreshTokenRequest) error
	AssignRole(ctx context.Context, userId int, roleId int) error
	RemoveRole(ctx context.Context, userId int, roleId int) error
	InitClientCredentialsReset(ctx context.Context, req InitClientCredentialsResetRequest) error
}

type AuthService struct {
	authProvider  port.Provider
	emailProvider email.Provider
	roleService   role.Provider
}

func NewAuthService(
	AP port.Provider,
	EP email.Provider,
	RP role.Provider) Provider {
	return &AuthService{
		authProvider:  AP,
		emailProvider: EP,
		roleService:   RP}
}

func (a AuthService) DecodeToken(ctx context.Context, token string) (DecodeResult, error) {
	res, err := a.authProvider.DecodeToken(ctx, token)
	if err != nil {
		return DecodeResult{}, ErrFailedToDecodeToken
	}
	return DecodeResult{
		Claims: Claims(res.Claims),
	}, nil
}

func (a AuthService) RetrospectToken(ctx context.Context, token string) (RetrospectionResult, error) {
	result, err := a.authProvider.RetrospectToken(ctx, token)
	if err != nil {
		return RetrospectionResult{}, ErrUnknown
	}

	return RetrospectionResult{Active: *result.Active}, nil
}

func (a AuthService) ResetClientCredentials(ctx context.Context, req *ResetCredentialsRequest) error {
	if req.ResetToken == "" {
		return ErrTokenNotSupplied
	}

	claims, err := a.validateToken(req.ResetToken)
	if err != nil {
		return err
	}

	if exp, ok := claims["exp"].(float64); ok {
		if int64(exp) < time.Now().Unix() {
			return ErrTokenHasExpired
		}
	} else {
		log.Print("expiration claim not found")
		return ErrUnknown
	}

	userId, ok := claims["userId"].(string)
	if !ok {
		log.Print("invalid user ID claim")
		return ErrUnknown
	}

	// Reset password
	err = a.authProvider.ResetPassword(ctx, userId, req.NewPassword)
	if err != nil {
		log.Print("failed to reset client credentials:", err)
		return ErrUnknown
	}

	return nil
}

func (a *AuthService) validateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Print("unexpected signing method")
			return nil, ErrUnknown
		}
		return []byte("s3cureR@nd0mK3y1234567890!"), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Print("invalid claims")
		return nil, ErrUnknown
	}

	return claims, nil
}

func (s AuthService) InitClientCredentialsReset(ctx context.Context, req InitClientCredentialsResetRequest) error {
	if req.UserId == "" {
		return ErrUserIdNotSupplied
	}

	if req.Email == "" {
		return ErrUserIdNotSupplied
	}

	_, err := s.createToken(req.UserId)
	if err != nil {
		return err
	}

	// Temp
	// err = s.sendResetEmail(token, req.Email)
	// if err != nil {
	// 	return err
	// }

	return nil
}

func (s AuthService) createToken(userId string) (string, error) {
	// Set token expiration time
	expirationTime := time.Now().Add(30 * time.Minute)
	claims := jwt.MapClaims{
		"exp":    expirationTime.Unix(),
		"userId": userId,
	}

	// Create a new token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token using a secret key
	signedToken, err := token.SignedString([]byte("s3cureR@nd0mK3y1234567890!"))
	if err != nil {
		return "", err // Return error if signing fails
	}

	log.Printf("Tokens %v,", signedToken)

	return signedToken, nil
}

func (s AuthService) sendResetEmail(token, recepientEmail string) error {
	err := s.emailProvider.Send(&email.SendRequest{
		To:         recepientEmail,
		Subject:    "Welcome to EfoytaStore!",
		TemplateId: 1,
		Args: map[string]string{
			"username":   recepientEmail,
			"reset_link": token,
		}})
	if err != nil {
		return ErrUnknown
	}
	return nil
}

func (a AuthService) CreateNewClientWithPassword(ctx context.Context, req *RegisterUserRequest) (RegisterUserResponse, error) {
	err := validateName(req.FirstName, req.LastName)
	if err != nil {
		return RegisterUserResponse{}, err
	}
	err = validatePasswords(req.Password)
	if err != nil {
		return RegisterUserResponse{}, err
	}
	err = validateEmail(req.Email)
	if err != nil {
		return RegisterUserResponse{}, err
	}
	err = validateUsername(req.Username)
	if err != nil {
		return RegisterUserResponse{}, err
	}

	resp, err := a.authProvider.CreateNewClient(ctx, &port.RegisterUserRequest{
		Email:       req.Email,
		Password:    req.Password,
		BirthDate:   req.BirthDate,
		PhoneNumber: req.PhoneNumber,
		ExternalId:  req.ExternalId,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Username:    req.Username})
	if err != nil {
		switch err {
		case port.ErrSysUsernameTaken:
			return RegisterUserResponse{}, ErrUsernameTaken
		case port.ErrSysEmailTaken:
			return RegisterUserResponse{}, ErrEmailTaken
		default:
			return RegisterUserResponse{}, ErrUnknown
		}
	}
	return RegisterUserResponse{
		Id:       resp.Id,
		Username: resp.Username,
	}, nil
}

func (a AuthService) CreateNewClientWithOutPassword(ctx context.Context, req *RegisterUserWithoutPasswordRequest) (RegisterUserResponse, error) {
	err := validateName(req.FirstName, req.LastName)
	if err != nil {
		return RegisterUserResponse{}, err
	}
	err = validateEmail(req.Email)
	if err != nil {
		return RegisterUserResponse{}, err
	}
	err = validateUsername(req.Username)
	if err != nil {
		return RegisterUserResponse{}, nil
	}

	genPassword, err := generateRandomPassword(10)
	if err != nil {
		log.Print("failed to generate password")
		return RegisterUserResponse{}, ErrUnknown
	}
	log.Printf("generated password: %v", genPassword)
	resp, err := a.authProvider.CreateNewClient(ctx, &port.RegisterUserRequest{
		Email:       req.Email,
		Password:    genPassword,
		BirthDate:   req.BirthDate,
		PhoneNumber: req.PhoneNumber,
		ExternalId:  req.ExternalId,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Username:    req.Username,
	})
	if err != nil {
		switch err {
		case port.ErrSysUsernameTaken:
			return RegisterUserResponse{}, ErrUsernameTaken
		case port.ErrSysEmailTaken:
			return RegisterUserResponse{}, ErrEmailTaken
		default:
			return RegisterUserResponse{}, ErrUnknown
		}
	}
	// err = a.InitClientCredentialsReset(ctx, InitClientCredentialsResetRequest{
	// 	UserId: resp.Id,
	// 	Email:  req.Email,
	// })
	// if err != nil {
	// 	log.Printf("Failed to init client credentials reset")
	// 	// a.authProvider.DeleteClient(ctx, req)
	// 	return RegisterUserResponse{}, err
	// }

	return RegisterUserResponse{
		Id:       resp.Id,
		Username: resp.Username,
	}, nil
}

func generateRandomPassword(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()"
	password := make([]byte, length)

	for i := 0; i < length; i++ {
		randIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		password[i] = charset[randIndex.Int64()]
	}

	return string(password), nil
}

func (a AuthService) DeleteClient(ctx context.Context, userId string) error {
	err := a.authProvider.DeleteClient(ctx, userId)
	if err != nil {
		return ErrUnknown
	}
	return nil
}

func (a AuthService) ClientLogin(ctx context.Context, req *LoginUserRequest) (LoginAuthResponse, error) {
	resp, err := a.authProvider.ClientLogin(ctx, &port.LoginUserRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		switch err {
		case port.ErrSysFailedToLogin:
			return LoginAuthResponse{}, ErrFailedToLogin
		default:
			return LoginAuthResponse{}, ErrUnknown
		}
	}
	return LoginAuthResponse{
		JWT: JWT(resp.JWT),
	}, nil
}

func (a AuthService) RefreshToken(ctx context.Context, req *RefreshTokenRequest) (LoginAuthResponse, error) {
	resp, err := a.authProvider.RefreshToken(ctx, &port.RefreshTokenRequest{RefreshToken: req.RefreshToken})
	if err != nil {
		switch err {
		case port.ErrSysFailedToLogin:
			return LoginAuthResponse{}, ErrFailedToLogin
		default:
			return LoginAuthResponse{}, ErrUnknown
		}
	}
	return LoginAuthResponse{
		JWT: JWT(resp.JWT),
	}, nil
}

func (a *AuthService) ClientLogout(ctx context.Context, req RefreshTokenRequest) error {
	err := a.authProvider.ClientLogout(ctx, port.RefreshTokenRequest(req))
	if err != nil {
		switch err {
		case port.ErrSysFailedToLogin:
			return ErrFailedToLogin
		default:
			return ErrUnknown
		}
	}
	return nil
}
func (a *AuthService) AssignRole(ctx context.Context, userId int, roleId int) error {
	return nil
}

func (a *AuthService) RemoveRole(ctx context.Context, userId int, roleId int) error {
	return nil
}

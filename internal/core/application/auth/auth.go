package auth

import (
	"context"
	"errors"
	"log"
	"time"

	"b2b.nati011.github.com/internal/core/application/email"
	port "b2b.nati011.github.com/internal/port/application/auth/provider"
	"github.com/golang-jwt/jwt"
)

var (
	ErrUsernameTaken        = errors.New("oopsy, username already taken")
	ErrEmailTaken           = errors.New("oopsy, email already taken")
	ErrFailedToLogin        = errors.New("oopsy, email or password incorrect")
	ErrUnknown              = errors.New("oopsy, unknown error has occured")
	ErrEmailNotSupplied     = errors.New("oopsy, email mandatory")
	ErrUsernameNotSupplied  = errors.New("oopsy, username mandatory")
	ErrInvalidEmail         = errors.New("oopsy, Invalid Email")
	ErrPasswordNotSupplied  = errors.New("oopsy, password mandatory")
	ErrFirstNameNotSupplied = errors.New("oopsy, First Name mandatory")
	ErrLastNameNotSupplied  = errors.New("oppsy, Last Name mandatory")
	ErrTokenNotSupplied     = errors.New("oppsy, Token is mandatory")
	ErrTokenHasExpired      = errors.New("oppsy, Token has expired")
	ErrUserIdNotSupplied    = errors.New("oopsy, userId mandatory")
)

type RegisterUserRequest struct {
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	BirthDate   time.Time `json:"birth_date"`
	PhoneNumber string    `json:"phone_number"`
	ExternalId  string    `json:"external_id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Username    string    `json:"username"`
}

type RegisterUserResponse struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type LoginAuthResponse struct {
	JWT JWT `json:"jwt"`
}

type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ResetCredentialsRequest struct {
	NewPassword string `json:"password"`
	ResetToken  string `json:"resetToken"`
}

type InitClientCredentialsResetRequest struct {
	UserId string
	Email  string
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

type Provider interface {
	CreateNewClient(ctx context.Context, req RegisterUserRequest) (RegisterUserResponse, error)
	ClientLogin(ctx context.Context, req LoginUserRequest) (LoginAuthResponse, error)
	RefreshToken(ctx context.Context, req RefreshTokenRequest) (LoginAuthResponse, error)
	DeleteClient(ctx context.Context, userId string) error
	ResetClientCredentials(ctx context.Context, req ResetCredentialsRequest) error
}

type AuthService struct {
	authProvider  port.Provider
	emailProvider email.Provider
}

func NewAuthService(ap port.Provider, em email.Provider) Provider {
	return &AuthService{authProvider: ap, emailProvider: em}
}

func (a *AuthService) ResetClientCredentials(ctx context.Context, req ResetCredentialsRequest) error {
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
		return errors.New("expiration claim not found")
	}

	userId, ok := claims["userId"].(string)
	if !ok {
		return errors.New("invalid user ID claim")
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

	token, err := s.createToken(req.UserId, req.Email)
	if err != nil {
		return err
	}

	err = s.sendResetEmail(token, req.Email)
	if err != nil {
		return err
	}

	return nil
}

func (s AuthService) createToken(userId, userEmail string) (string, error) {
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

func (a *AuthService) CreateNewClient(ctx context.Context, req RegisterUserRequest) (RegisterUserResponse, error) {
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

	resp, err := a.authProvider.CreateNewClient(ctx, port.RegisterUserRequest(req))
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
	err = a.InitClientCredentialsReset(ctx, InitClientCredentialsResetRequest{
		UserId: resp.Id,
		Email:  req.Email,
	})
	if err != nil {
		log.Printf("Failed to init client credentials reset")
		// a.authProvider.DeleteClient(ctx, req)
		return RegisterUserResponse{}, err
	}

	return RegisterUserResponse{
		Id:       resp.Id,
		Username: resp.Username,
	}, nil
}

func (a *AuthService) DeleteClient(ctx context.Context, userId string) error {
	err := a.authProvider.DeleteClient(ctx, userId)
	if err != nil {
		return ErrUnknown
	}
	return nil
}

func (a *AuthService) ClientLogin(ctx context.Context, rq LoginUserRequest) (LoginAuthResponse, error) {
	resp, err := a.authProvider.ClientLogin(ctx, port.LoginUserRequest(rq))
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
func (a *AuthService) RefreshToken(ctx context.Context, req RefreshTokenRequest) (LoginAuthResponse, error) {
	resp, err := a.authProvider.RefreshToken(ctx, port.RefreshTokenRequest(req))
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

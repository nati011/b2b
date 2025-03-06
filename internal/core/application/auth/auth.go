package auth

import (
	"errors"

	port "b2b.nati011.github.com/internal/port/application/auth/provider"
)

// user readable errors
var (
	ErrUsernameTaken = errors.New("oopsy, username already taken")
	ErrEmailTaken    = errors.New("oopsy, email already taken")
	ErrFailedToLogin = errors.New("oopsy, email or password incorrect")
	ErrUnknown       = errors.New("oopsy, unknown error has occured")
)

type RegisterUserRequest struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmed_password"`
	FullName        string `json:"full_name"`
	Username        string `json:"username"`
}

type RegisterUserResponse struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUserResonse struct {
	JWT     JWT
	Message string
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

const (
	SUCCESS_MESSAGE = "Ahoy!"
)

type Provider interface {
	CreateClient(RegisterUserRequest) (RegisterUserResponse, error)
	LoginClient(LoginUserRequest) (LoginUserResonse, error)
}

type AuthService struct {
	authProvider port.AuthProvider
}

func NewAuthService(ap port.AuthProvider) *AuthService {
	return &AuthService{authProvider: ap}
}

func (a AuthService) CreateClient(rq RegisterUserRequest) (RegisterUserResponse, error) {
	resp, err := a.authProvider.CreateNewClient(rq.FullName, rq.FullName, rq.Email, rq.Username, rq.Password)
	if err != nil {
		switch err {
		case port.ErrSysUsernameTaken:
			return RegisterUserResponse{
				Message: ErrUsernameTaken.Error(),
			}, ErrUsernameTaken
		case port.ErrSysEmailTaken:
			return RegisterUserResponse{
				Message: ErrEmailTaken.Error(),
			}, ErrEmailTaken
		default:
			return RegisterUserResponse{}, ErrUnknown
		}
	}
	return RegisterUserResponse{
		Username: resp.Username,
		Message:  SUCCESS_MESSAGE,
	}, nil
}

func (a *AuthService) LoginClient(rq LoginUserRequest) (LoginUserResonse, error) {
	resp, err := a.authProvider.ClientLogin(rq.Email, rq.Password)
	if err != nil {
		switch err {
		case port.ErrSysFailedToLogin:
			return LoginUserResonse{
				Message: ErrFailedToLogin.Error(),
			}, ErrFailedToLogin
		default:
			return LoginUserResonse{}, ErrUnknown
		}
	}
	return LoginUserResonse{
		JWT:     JWT(resp.JWT),
		Message: SUCCESS_MESSAGE,
	}, nil
}

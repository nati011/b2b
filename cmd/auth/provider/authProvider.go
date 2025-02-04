package auth

import "errors"

// system errors
var (
	System_Readable_Error_UsernameTaken = errors.New("invalid username")
	System_Readable_Error_EmailTaken    = errors.New("invalid email")
	System_Readable_Error_FailedToLogin = errors.New("invalid email or password")
)

type LoginAuthResonse struct {
	JWT struct {
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
}

type CreateClientAuthResonse struct {
	Username string
}

type AuthProvider interface {
	CreateNewClient(firstName string, lastName string, email string, username string, password string) (CreateClientAuthResonse, error)
	ClientLogin(email, password string) (LoginAuthResonse, error)
}

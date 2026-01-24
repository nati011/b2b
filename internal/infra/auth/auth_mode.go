package auth

// AuthMode represents the authentication mode
type AuthMode string

const (
	AuthModeBasic  AuthMode = "basic"
	AuthModeOAuth2 AuthMode = "oauth2"
)

// String returns the string representation of AuthMode
func (a AuthMode) String() string {
	return string(a)
}

// IsValid checks if the auth mode is valid
func (a AuthMode) IsValid() bool {
	return a == AuthModeBasic || a == AuthModeOAuth2
}

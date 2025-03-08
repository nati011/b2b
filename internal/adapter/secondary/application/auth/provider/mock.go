package provider

import (
	port "b2b.nati011.github.com/internal/port/application/auth/provider"
)

type MockClient struct {
	firstName string
	lastName  string
	email     string
	username  string
	password  string
}

type MockAuthProvider struct {
	clients     []MockClient
	validTokens []string
}

// RefreshToken implements provider.AuthProvider.
func (m *MockAuthProvider) RefreshToken(token string) (port.LoginAuthResponse, error) {
	panic("unimplemented")
}

func NewMockAuthProvider() MockAuthProvider {
	return MockAuthProvider{}
}

func (m *MockAuthProvider) Cleanup() {
	m.clients = []MockClient{}
}

func (m *MockAuthProvider) CreateNewClient(firstName string, lastName string, email string, username string, password string) (port.CreateClientAuthResponse, error) {
	// check if username or password is taken
	for _, index := range m.clients {
		if index.username == username {
			return port.CreateClientAuthResponse{}, port.ErrSysUsernameTaken
		} else if index.email == email {
			return port.CreateClientAuthResponse{}, port.ErrSysEmailTaken
		}
	}
	m.clients = append(m.clients, MockClient{
		firstName: firstName,
		lastName:  lastName,
		email:     email,
		username:  username,
		password:  password,
	})

	return port.CreateClientAuthResponse{
		Username: username,
	}, nil
}

func (m *MockAuthProvider) ClientLogin(email, password string) (port.LoginAuthResponse, error) {
	// check if username or password is taken
	clientExists := false

	for _, index := range m.clients {
		if index.email == email {
			clientExists = true
		}
	}
	if !clientExists {
		return port.LoginAuthResponse{}, port.ErrSysFailedToLogin
	}

	return port.LoginAuthResponse{}, nil
}

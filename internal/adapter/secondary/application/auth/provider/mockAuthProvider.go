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
	clients []MockClient
}

func NewMockAuthProvider() port.AuthProvider {
	return &MockAuthProvider{}
}

func (m *MockAuthProvider) FlushMockAuthProvider() {
	m.clients = []MockClient{}
}

func (m *MockAuthProvider) CreateNewClient(firstName string, lastName string, email string, username string, password string) (port.CreateClientAuthResonse, error) {
	// check if username or password is taken
	for _, index := range m.clients {
		if index.username == username {
			return port.CreateClientAuthResonse{}, port.ErrSysUsernameTaken
		} else if index.email == email {
			return port.CreateClientAuthResonse{}, port.ErrSysEmailTaken
		}
	}
	m.clients = append(m.clients, MockClient{
		firstName: firstName,
		lastName:  lastName,
		email:     email,
		username:  username,
		password:  password,
	})

	return port.CreateClientAuthResonse{
		Username: username,
	}, nil
}

func (m *MockAuthProvider) ClientLogin(email, password string) (port.LoginAuthResonse, error) {
	// check if username or password is taken
	clientExists := false

	for _, index := range m.clients {
		if index.email == email {
			clientExists = true
		}
	}
	if !clientExists {
		return port.LoginAuthResonse{}, port.ErrSysFailedToLogin
	}

	return port.LoginAuthResonse{}, nil
}

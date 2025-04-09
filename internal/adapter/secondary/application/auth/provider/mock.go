package provider

import (
	"context"

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

func (m *MockAuthProvider) RefreshToken(ctx context.Context, req port.RefreshTokenRequest) (port.LoginAuthResponse, error) {
	panic("unimplemented")
}

func NewMockAuthProvider() MockAuthProvider {
	return MockAuthProvider{}
}

func (m *MockAuthProvider) Teardown() {
	m.clients = []MockClient{}
}

func (m *MockAuthProvider) DeleteClient(ctx context.Context, userId string) error {
	panic("unimplemented")
}

func (m *MockAuthProvider) CreateNewClient(ctx context.Context, req port.RegisterUserRequest) (port.RegisterUserResponse, error) {
	// check if username or password is taken
	for _, index := range m.clients {
		if index.username == req.Username {
			return port.RegisterUserResponse{}, port.ErrSysUsernameTaken
		} else if index.email == req.Email {
			return port.RegisterUserResponse{}, port.ErrSysEmailTaken
		}
	}
	m.clients = append(m.clients, MockClient{
		firstName: req.FirstName,
		lastName:  req.LastName,
		email:     req.Email,
		username:  req.Username,
		password:  req.Password,
	})

	return port.RegisterUserResponse{
		Username: req.Username,
	}, nil
}

func (m *MockAuthProvider) ClientLogin(ctx context.Context, req port.LoginUserRequest) (port.LoginAuthResponse, error) {
	// check if username or password is taken
	clientExists := false

	for _, index := range m.clients {
		if index.email == req.Email {
			clientExists = true
		}
	}
	if !clientExists {
		return port.LoginAuthResponse{}, port.ErrSysFailedToLogin
	}

	return port.LoginAuthResponse{}, nil
}

package provider

import (
	"context"
	"strconv"

	port "b2b.nati011.github.com/internal/port/application/auth/provider"
)

type MockClient struct {
	userId    string
	firstName string
	lastName  string
	email     string
	username  string
	password  string
	roles     []int
}

type MockAuthProvider struct {
	clients []MockClient
}

func NewMockAuthProvider() MockAuthProvider {
	return MockAuthProvider{}
}

func (m MockAuthProvider) Teardown() {
	m.clients = []MockClient{}
}

func (m MockAuthProvider) RefreshToken(ctx context.Context, req *port.RefreshTokenRequest) (port.LoginAuthResponse, error) {
	return port.LoginAuthResponse{}, nil
}

func (m MockAuthProvider) GoogleSignOn(ctx context.Context, SubjectToken string) (port.LoginAuthResponse, error) {
	return port.LoginAuthResponse{}, nil
}

func (m MockAuthProvider) RetrospectToken(ctx context.Context, token string) (port.RetrospectionResult, error) {
	return port.RetrospectionResult{}, nil
}

func (m MockAuthProvider) DecodeToken(ctx context.Context, token string) (port.DecodedResult, error) {
	return port.DecodedResult{}, nil
}

func (m MockAuthProvider) DeleteClient(ctx context.Context, userId string) error {
	for _, i := range m.clients {
		if i.userId != userId {
			m.clients = append(m.clients, MockClient{
				userId:    i.userId,
				firstName: i.firstName,
				lastName:  i.lastName,
				email:     i.email,
				username:  i.username,
				password:  i.password,
			})
		}
	}
	return nil
}

func (m MockAuthProvider) CreateNewClient(ctx context.Context, req *port.RegisterUserRequest) (port.RegisterUserResponse, error) {
	newId := strconv.Itoa(len(m.clients) + 1)
	// check if username or email is taken
	for _, index := range m.clients {
		if index.username == req.Username {
			return port.RegisterUserResponse{}, port.ErrSysUsernameTaken
		} else if index.email == req.Email {
			return port.RegisterUserResponse{}, port.ErrSysEmailTaken
		}
	}

	m.clients = append(m.clients, MockClient{
		userId:    newId,
		firstName: req.FirstName,
		lastName:  req.LastName,
		email:     req.Email,
		username:  req.Username,
		password:  req.Password,
	})

	return port.RegisterUserResponse{
		Id:       newId,
		Username: req.Username,
	}, nil
}

func (m MockAuthProvider) ClientLogin(ctx context.Context, req *port.LoginUserRequest) (port.LoginAuthResponse, error) {
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

func (k MockAuthProvider) ClientLogout(ctx context.Context, req port.RefreshTokenRequest) error {
	return nil
}

func (k MockAuthProvider) ResetPassword(ctx context.Context, userId, new_password string) error {
	updated := []MockClient{}
	for _, i := range k.clients {
		if i.userId == userId {
			updated = append(updated, MockClient{
				userId:    i.userId,
				firstName: i.firstName,
				lastName:  i.lastName,
				email:     i.email,
				username:  i.username,
				password:  new_password,
			})
		} else {
			updated = append(updated, MockClient{
				userId:    i.userId,
				firstName: i.firstName,
				lastName:  i.lastName,
				email:     i.email,
				username:  i.username,
				password:  i.password,
			})
		}
	}
	k.clients = updated
	return nil
}

func (m MockAuthProvider) AssignRole(ctx context.Context, userId string, roleId int) error {
	for _, i := range m.clients {
		if i.userId == userId {
			i.roles = append(i.roles, roleId)
		}
	}
	return nil
}

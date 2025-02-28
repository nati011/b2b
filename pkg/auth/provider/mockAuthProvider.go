package provider

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

func NewMockAuthProvider() AuthProvider {
	return &MockAuthProvider{}
}

func (m *MockAuthProvider) FlushMockAuthProvider() {
	m.clients = []MockClient{}
}

func (m *MockAuthProvider) CreateNewClient(firstName string, lastName string, email string, username string, password string) (CreateClientAuthResponse, error) {
	// check if username or password is taken
	for _, index := range m.clients {
		if index.username == username {
			return CreateClientAuthResponse{}, ErrSysUsernameTaken
		} else if index.email == email {
			return CreateClientAuthResponse{}, ErrSysEmailTaken
		}
	}
	m.clients = append(m.clients, MockClient{
		firstName: firstName,
		lastName:  lastName,
		email:     email,
		username:  username,
		password:  password,
	})

	return CreateClientAuthResponse{
		Username: username,
	}, nil
}

func (m *MockAuthProvider) ClientLogin(email, password string) (LoginAuthResponse, error) {
	// check if username or password is taken
	clientExists := false

	for _, index := range m.clients {
		if index.email == email {
			clientExists = true
		}
	}
	if !clientExists {
		return LoginAuthResponse{}, ErrSysFailedToLogin
	}

	return LoginAuthResponse{}, nil
}

func (m *MockAuthProvider) RefreshToken(token string) (LoginAuthResponse, error) {
	validToken := false
	for _, index := range m.validTokens {
		if index == token {
			validToken = true
		}
	}
	if !validToken {
		return LoginAuthResponse{}, ErrSysTokenExpired
	}

	return LoginAuthResponse{}, nil
}

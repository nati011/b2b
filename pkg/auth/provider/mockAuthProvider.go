package provider

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

func NewMockAuthProvider() AuthProvider {
	return &MockAuthProvider{}
}

func (m *MockAuthProvider) FlushMockAuthProvider() {
	m.clients = []MockClient{}
}

func (m *MockAuthProvider) CreateNewClient(firstName string, lastName string, email string, username string, password string) (CreateClientAuthResonse, error) {
	// check if username or password is taken
	for _, index := range m.clients {
		if index.username == username {
			return CreateClientAuthResonse{}, ErrSysUsernameTaken
		} else if index.email == email {
			return CreateClientAuthResonse{}, ErrSysEmailTaken
		}
	}
	m.clients = append(m.clients, MockClient{
		firstName: firstName,
		lastName:  lastName,
		email:     email,
		username:  username,
		password:  password,
	})

	return CreateClientAuthResonse{
		Username: username,
	}, nil
}

func (m *MockAuthProvider) ClientLogin(email, password string) (LoginAuthResonse, error) {
	// check if username or password is taken
	clientExists := false

	for _, index := range m.clients {
		if index.email == email {
			clientExists = true
		}
	}
	if !clientExists {
		return LoginAuthResonse{}, ErrSysFailedToLogin
	}

	return LoginAuthResonse{}, nil
}

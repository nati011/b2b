package auth

import (
	"testing"

	authDTO "b2b.nati011.github.com/cmd/auth/model/dto"
	authProvider "b2b.nati011.github.com/cmd/auth/provider"
)

const (
	VALID_PASSWORD = "test@123"
	VALID_FULLNAME = "ruth tirusew"

	//VALID
	VALID_USERNAME_a = "expired_pineapple"
	VALID_EMAIL_a    = "ruthtirusew944@gmail.com"

	VALID_USERNAME_b = "expired_pineapple"
	VALID_EMAIL_b    = "ruthtirusew944@gmail.com"

	//INVALID
	INVALID_username = ""
	INVALID_PASSWORD = ""
	INVALID_FULLNAME = ""
	INVALID_EMAIL    = ""
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

func NewMockAuthProvider() authProvider.AuthProvider {
	return &MockAuthProvider{}
}

func (m *MockAuthProvider) CreateNewClient(firstName string, lastName string, email string, username string, password string) (authProvider.CreateClientAuthResonse, error) {

	// check if username or password is taken
	for _, index := range m.clients {
		if index.username == username {
			return authProvider.CreateClientAuthResonse{}, ErrUsernameTaken
		} else if index.email == email {
			return authProvider.CreateClientAuthResonse{}, ErrEmailTaken
		}
	}
	m.clients = append(m.clients, MockClient{
		firstName: firstName,
		lastName:  lastName,
		email:     email,
		username:  username,
		password:  password,
	})

	return authProvider.CreateClientAuthResonse{}, nil
}

func (m *MockAuthProvider) ClientLogin(email, password string) (authProvider.LoginAuthResonse, error) {

	// check if username or password is taken
	for _, index := range m.clients {
		if index.email == email {
			return authProvider.LoginAuthResonse{}, ErrUsernameTaken
		} else {
			return authProvider.LoginAuthResonse{}, ErrEmailTaken
		}
	}

	return authProvider.LoginAuthResonse{}, nil
}

var mockProvider = NewMockAuthProvider()
var container = NewContainer(mockProvider)

func TestCreateClient_happyPath(t *testing.T) {

	user := authDTO.RegisterUserRequest{
		Username:        VALID_USERNAME_a,
		Password:        VALID_PASSWORD,
		ConfirmPassword: VALID_PASSWORD,
		FullName:        VALID_FULLNAME,
		Email:           VALID_EMAIL_a,
	}

	userRegistrationSuccessResponse := authDTO.RegisterUserResponse{
		Username: VALID_USERNAME_a,
		Message:  SUCCESS_MESSAGE,
	}

	in := user
	want := userRegistrationSuccessResponse

	got, err := container.AuthService.CreateClient(in)
	if err != nil {
		t.Errorf("failed to create client err: %v", err)
	}
	if got != want {
		t.Errorf("want:%v \n got:%v", want, got)
	}
}

func TestCreateClient_unhappyPath(t *testing.T) {

	t.Run("duplicateUsername", func(t *testing.T) {

		//init
		in_a := authDTO.RegisterUserRequest{
			Username:        VALID_USERNAME_a,
			Password:        VALID_PASSWORD,
			ConfirmPassword: VALID_PASSWORD,
			FullName:        VALID_FULLNAME,
			Email:           VALID_EMAIL_a,
		}

		got, err := container.AuthService.CreateClient(in_a)
		if err != nil {
			t.Errorf("failed to create client err: %v", err)
		}

		in_b := authDTO.RegisterUserRequest{
			Username:        VALID_USERNAME_a,
			Password:        VALID_PASSWORD,
			ConfirmPassword: VALID_PASSWORD,
			FullName:        VALID_FULLNAME,
			Email:           VALID_EMAIL_b,
		}

		want := authDTO.RegisterUserResponse{
			Username: "",
			Message:  ErrUsernameTaken.Error(),
		}

		got, err = container.AuthService.CreateClient(in_b)
		if err != nil {
			t.Errorf("failed to create client err: %v", err)
		}

		if got != want {
			t.Errorf("want:%v \n got:%v", want, got)
		}

	})

	t.Run("duplicateEmail", func(t *testing.T) {

		//init
		ua := authDTO.RegisterUserRequest{
			Username:        "expired_pineapple",
			Password:        "test@123",
			ConfirmPassword: "test@123",
			FullName:        "ruth tirusew",
			Email:           "ruthtirusew944@gmail.com",
		}

		container.AuthService.CreateClient(ua)
		ub := authDTO.RegisterUserRequest{
			Username:        "expired_pineapple",
			Password:        "test@123",
			ConfirmPassword: "test@123",
			FullName:        "ruth tirusew",
			Email:           "ruthtirusew30@gmail.com",
		}

		want := authDTO.RegisterUserResponse{
			Username: "",
			Message:  ErrEmailTaken.Error(),
		}

		got, err := container.AuthService.CreateClient(ub)
		if err != nil {
			t.Errorf("failed to create client err: %v", err)
		}

		if got != want {
			t.Errorf("want:%v \n got:%v", want, got)
		}

	})
}

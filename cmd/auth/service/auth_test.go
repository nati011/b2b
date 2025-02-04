package auth

import (
	"testing"

	"b2b.nati011.github.com/cmd/auth"
	authDTO "b2b.nati011.github.com/cmd/auth/model/dto"
	authProvider "b2b.nati011.github.com/cmd/auth/provider"
)

const (
	//VALID CREDENTIALS
	VALID_USERNAME_a = "expired_pineapple"
	VALID_PASSWORD_a = "test@123"
	VALID_FULLNAME_a = "ruth tirusew"
	VALID_EMAIL_a    = "ruthtirusew944@gmail.com"

	VALID_USERNAME_b = "expired_pineapple"
	VALID_PASSWORD_b = "test@123"
	VALID_FULLNAME_b = "ruth tirusew"
	VALID_EMAIL_b    = "ruthtirusew944@gmail.com"

	//INVALID CREDENTIALS
	INVALID_username = ""
	INVALID_PASSWORD = ""
	INVALID_FULLNAME = ""
	INVALID_EMAIL    = ""

	//MESSAGES
	SUCCESS_REGISTRATION_MESSAGE                  = "Ahoy"
	FAILURE_REGISTRATION_MESSAGE_INVALID_USERNAME = "Oopsy, username is already taken"
	FAILURE_REGISTRATION_MESSAGE_INVALID_EMAIL    = "Oopsy, email is already taken"

	// ??SHOULD BE MESSAGES OR ERROR(PERHAPS HUMAN READABLE ERRORS)??
	SUCCESS_SIGNIN_MESSAGE                           = "Ahoy"
	FAILURE_REGISTRATION_MESSAGE_INVALID_CREDENTAILS = "Oopsy, email or password incorrect"

	//ERRORS
)

// const authService = NewAuthService(mockProvider)

type MockAuthProvider struct {
}

func NewMockAuthProvider() authProvider.AuthProvider {
	return &MockAuthProvider{}
}

func (m *MockAuthProvider) CreateNewClient(firstName string, lastName string, email string, username string, password string) (authProvider.CreateClientAuthResonse, error) {
	type Client struct {
		firstName string
		lastName  string
		email     string
		username  string
		password  string
	}

	users := []Client{}

	// check if username exists
	for _, index := range users {
		if index.username == username {
			return ClientUsernameAlreadyUsed(username)
		} else if index.email == email {
			return ClientEmailAlreadyUsed(email)
		}
	}
	users = append(users, Client{
		firstName: firstName,
		lastName:  lastName,
		email:     email,
		username:  username,
		password:  password,
	})

	return authProvider.CreateClientAuthResonse{}, nil
}

func (m *MockAuthProvider) ClientLogin(email, password string) (authProvider.LoginAuthResonse, error) {

	return authProvider.LoginAuthResonse{}, nil
}

var mockProvider = NewMockAuthProvider()
var container = auth.NewContainer(mockProvider)

func TestCreateClient_happyPath(t *testing.T) {

	user := authDTO.RegisterUserRequest{
		Username:        VALID_USERNAME_a,
		Password:        "test@123",
		ConfirmPassword: "test@123",
		FullName:        "ruth tirusew",
		Email:           "ruthtirusew944@gmail.com",
	}

	userRegistrationSuccessResponse := authDTO.RegisterUserResponse{
		Username: VALID_USERNAME_a,
		Message:  SUCCESS_REGISTRATION_MESSAGE,
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
			Username:        "expired_pineapple",
			Password:        "test@123",
			ConfirmPassword: "test@123",
			FullName:        "ruth tirusew",
			Email:           "ruthtirusew944@gmail.com",
		}

		got, err := container.AuthService.CreateClient(in_a)
		if err != nil {
			t.Errorf("failed to create client err: %v", err)
		}

		in_b := authDTO.RegisterUserRequest{
			Username:        "expired_pineapple",
			Password:        "test@123",
			ConfirmPassword: "test@123",
			FullName:        "ruth tirusew",
			Email:           "ruthtirusew30@gmail.com",
		}

		want := authDTO.RegisterUserResponse{
			Username: "",
			Message:  FAILURE_REGISTRATION_MESSAGE_INVALID_USERNAME,
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
			Message:  FAILURE_REGISTRATION_MESSAGE_INVALID_EMAIL,
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

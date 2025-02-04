package auth

import (
	"testing"

	authDTO "b2b.nati011.github.com/cmd/auth/model/dto"
	"b2b.nati011.github.com/cmd/auth/provider"
)

const (
	VALID_PASSWORD = "test@123"
	VALID_FULLNAME = "ruth tirusew"
	VALID_EMAIL    = "ruthtirusew944@gmail.com"
	VALID_USERNAME = "expired_pineapple"

	//INVALID
	INVALID_username = ""
	INVALID_PASSWORD = ""
	INVALID_FULLNAME = ""
	INVALID_EMAIL    = ""
)

var mockProvider = provider.NewMockAuthProvider()
var container = NewContainer(mockProvider)

func TestCreateClient_happyPath(t *testing.T) {

	user := authDTO.RegisterUserRequest{
		Username:        VALID_USERNAME,
		Password:        VALID_PASSWORD,
		ConfirmPassword: VALID_PASSWORD,
		FullName:        VALID_FULLNAME,
		Email:           VALID_EMAIL,
	}

	userRegistrationSuccessResponse := authDTO.RegisterUserResponse{
		Username: VALID_USERNAME,
		Message:  SUCCESS_MESSAGE,
	}

	in := user
	want := userRegistrationSuccessResponse

	got, err := container.AuthService.CreateClient(in)

	if err != nil {
		t.Errorf("Failed to create client err: %v", err)
	}
	if got != want {
		t.Errorf("Expected: %v, Got: %v", want, got)
	}
}

func TestCreateClient_UnhappyPath(t *testing.T) {

	t.Run("DuplicateUsername", func(t *testing.T) {
		//init
		in_a := authDTO.RegisterUserRequest{
			Username:        VALID_USERNAME,
			Password:        VALID_PASSWORD,
			ConfirmPassword: VALID_PASSWORD,
			FullName:        VALID_FULLNAME,
			Email:           VALID_EMAIL,
		}

		_, err := container.AuthService.CreateClient(in_a)
		if err != nil {
			t.Errorf("Failed to create client err: %v", err)
		}

		in_b := authDTO.RegisterUserRequest{
			Username:        VALID_USERNAME,
			Password:        VALID_PASSWORD,
			ConfirmPassword: VALID_PASSWORD,
			FullName:        VALID_FULLNAME,
			Email:           VALID_EMAIL,
		}

		want := authDTO.RegisterUserResponse{
			Username: "",
			Message:  ErrUsernameTaken.Error(),
		}

		got, err := container.AuthService.CreateClient(in_b)
		if err != nil {
			if got != want {
				t.Errorf("Expected: %v, Got: %v", want, got)
			}
		}

	})

	t.Run("DuplicateEmail", func(t *testing.T) {

		//init
		ua := authDTO.RegisterUserRequest{
			Username:        VALID_USERNAME,
			Password:        VALID_PASSWORD,
			ConfirmPassword: VALID_PASSWORD,
			FullName:        VALID_FULLNAME,
			Email:           VALID_EMAIL,
		}

		container.AuthService.CreateClient(ua)
		ub := authDTO.RegisterUserRequest{
			Username:        VALID_USERNAME,
			Password:        VALID_PASSWORD,
			ConfirmPassword: VALID_PASSWORD,
			FullName:        VALID_FULLNAME,
			Email:           VALID_EMAIL,
		}

		want := authDTO.RegisterUserResponse{
			Username: "",
			Message:  ErrEmailTaken.Error(),
		}

		got, err := container.AuthService.CreateClient(ub)
		if err != nil {
			if got != want {
				t.Errorf("Expected: %v, Got: %v", want, got)
			}
		}

	})
}

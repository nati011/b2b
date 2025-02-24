package auth

import (
	"testing"

	authDTO "b2b.nati011.github.com/internal/core/domain/auth/model/dto"
	"b2b.nati011.github.com/internal/core/domain/auth/provider"
)

const (
	VALID_PASSWORD   = "test@123"
	VALID_FULLNAME   = "ruth tirusew"
	VALID_EMAIL_A    = "ruthtirusew944@gmail.com"
	VALID_EMAIL_B    = "ruthtirusew388@gmail.com"
	VALID_USERNAME_A = "expired_pineapple"
	VALID_USERNAME_B = "delila"

	//INVALID
	INVALID_username = ""
	INVALID_PASSWORD = ""
	INVALID_FULLNAME = ""
	INVALID_EMAIL    = ""
)

var container = NewContainer(provider.NewMockAuthProvider())

func TestCreateClient_happyPath(t *testing.T) {
	user := authDTO.RegisterUserRequest{
		Username:        VALID_USERNAME_A,
		Password:        VALID_PASSWORD,
		ConfirmPassword: VALID_PASSWORD,
		FullName:        VALID_FULLNAME,
		Email:           VALID_EMAIL_A,
	}

	userRegistrationSuccessResponse := authDTO.RegisterUserResponse{
		Username: VALID_USERNAME_A,
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
			Username:        VALID_USERNAME_A,
			Password:        VALID_PASSWORD,
			ConfirmPassword: VALID_PASSWORD,
			FullName:        VALID_FULLNAME,
			Email:           VALID_EMAIL_A,
		}

		_, err := container.AuthService.CreateClient(in_a)
		if err != nil {
			t.Errorf("Failed to create client err: %v", err)
		}

		in_b := authDTO.RegisterUserRequest{
			Username:        VALID_USERNAME_A,
			Password:        VALID_PASSWORD,
			ConfirmPassword: VALID_PASSWORD,
			FullName:        VALID_FULLNAME,
			Email:           VALID_EMAIL_B,
		}

		want := authDTO.RegisterUserResponse{
			Username: "",
			Message:  ErrUsernameTaken.Error(),
		}

		got, _ := container.AuthService.CreateClient(in_b)
		if got != want {
			t.Errorf("Expected: %v, Got: %v", want, got)
		}

	})

	t.Run("DuplicateEmail", func(t *testing.T) {

		//init
		ua := authDTO.RegisterUserRequest{
			Username:        VALID_USERNAME_A,
			Password:        VALID_PASSWORD,
			ConfirmPassword: VALID_PASSWORD,
			FullName:        VALID_FULLNAME,
			Email:           VALID_EMAIL_A,
		}

		container.AuthService.CreateClient(ua)
		ub := authDTO.RegisterUserRequest{
			Username:        VALID_USERNAME_B,
			Password:        VALID_PASSWORD,
			ConfirmPassword: VALID_PASSWORD,
			FullName:        VALID_FULLNAME,
			Email:           VALID_EMAIL_A,
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

package auth

import (
	"os"
	"testing"

	"b2b.nati011.github.com/internal/adapter/secondary/application/auth/provider"
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

var service Provider
var mock = provider.NewMockAuthProvider()

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	service = NewAuthService(&mock)
}

func Test_CreateClient_happyPath(t *testing.T) {
	user := RegisterUserRequest{
		Username:        VALID_USERNAME_A,
		Password:        VALID_PASSWORD,
		ConfirmPassword: VALID_PASSWORD,
		FullName:        VALID_FULLNAME,
		Email:           VALID_EMAIL_A,
	}

	userRegistrationSuccessResponse := RegisterUserResponse{
		Username: VALID_USERNAME_A,
	}

	in := user
	want := userRegistrationSuccessResponse

	got, err := service.CreateClient(in)

	if err != nil {
		t.Fatalf("Failed to create client err: %v", err)
	}
	if got != want {
		t.Errorf("Expected: %v, Got: %v", want, got)
	}
	mock.Cleanup()
}

func Test_CreateClient_UnhappyPath(t *testing.T) {
	t.Run("DuplicateUsername", func(t *testing.T) {
		//setup
		in_a := RegisterUserRequest{
			Username:        VALID_USERNAME_A,
			Password:        VALID_PASSWORD,
			ConfirmPassword: VALID_PASSWORD,
			FullName:        VALID_FULLNAME,
			Email:           VALID_EMAIL_A,
		}

		_, err := service.CreateClient(in_a)
		if err != nil {
			t.Fatalf("Failed to create client err: %v", err)
		}

		in_b := RegisterUserRequest{
			Username:        VALID_USERNAME_A,
			Password:        VALID_PASSWORD,
			ConfirmPassword: VALID_PASSWORD,
			FullName:        VALID_FULLNAME,
			Email:           VALID_EMAIL_B,
		}

		want := RegisterUserResponse{
			Username: "",
		}

		got, _ := service.CreateClient(in_b)
		if got != want {
			t.Errorf("Expected: %v, Got: %v", want, got)
		}
		mock.Cleanup()
	})

	t.Run("DuplicateEmail", func(t *testing.T) {
		//setup
		ua := RegisterUserRequest{
			Username:        VALID_USERNAME_A,
			Password:        VALID_PASSWORD,
			ConfirmPassword: VALID_PASSWORD,
			FullName:        VALID_FULLNAME,
			Email:           VALID_EMAIL_A,
		}

		service.CreateClient(ua)
		ub := RegisterUserRequest{
			Username:        VALID_USERNAME_B,
			Password:        VALID_PASSWORD,
			ConfirmPassword: VALID_PASSWORD,
			FullName:        VALID_FULLNAME,
			Email:           VALID_EMAIL_A,
		}

		want := RegisterUserResponse{
			Username: "",
		}

		got, err := service.CreateClient(ub)
		if err != nil {
			if got != want {
				t.Errorf("Expected: %v, Got: %v", want, got)
			}
		}
		mock.Cleanup()
	})
}

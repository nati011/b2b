package auth

import (
	"context"
	"os"
	"testing"

	"b2b.nati011.github.com/internal/adapter/secondary/application/auth/provider"
)

const (
	VALID_PASSWORD   = "test@123"
	VALID_PASSWORD_B = "yesy"
	VALID_FirstName  = "tirusew"
	VALID_LastName   = "tirusew"
	VALID_EMAIL_A    = "ruthtirusew944@gmail.com"
	VALID_EMAIL_B    = "ruthtirusew388@gmail.com"
	VALID_USERNAME_A = "expired_pineapple"
	VALID_USERNAME_B = "delilah"

	//INVALID
	INVALID_username  = ""
	INVALID_PASSWORD  = ""
	INVALID_FirstName = ""
	INVALID_EMAIL     = ""
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
	t.Cleanup(mock.Teardown)
	ctx := context.Background()
	in := RegisterUserRequest{
		Username:  VALID_USERNAME_A,
		Password:  VALID_PASSWORD,
		FirstName: VALID_FirstName,
		LastName:  VALID_LastName,
		Email:     VALID_EMAIL_A,
	}
	_, err := service.CreateNewClient(ctx, in)
	if err != nil {
		t.Errorf("Failed to create err: %v", err)
	}
}

func Test_CreateClient_UnhappyPath(t *testing.T) {
	t.Run("email_not_supplied", func(t *testing.T) {
		t.Cleanup(mock.Teardown)
		ctx := context.Background()
		in := RegisterUserRequest{
			Username:  VALID_USERNAME_A,
			Password:  VALID_PASSWORD,
			FirstName: VALID_FirstName,
			LastName:  VALID_LastName,
		}
		_, err := service.CreateNewClient(ctx, in)
		wantErr := ErrEmailNotSupplied
		if err != wantErr {
			t.Errorf("Expected err: %v Got : %v", wantErr, err)
		}
	})

	t.Run("password_not_supplied", func(t *testing.T) {
		t.Cleanup(mock.Teardown)
		ctx := context.Background()
		in := RegisterUserRequest{
			Email:     VALID_EMAIL_A,
			Username:  VALID_USERNAME_A,
			FirstName: VALID_FirstName,
			LastName:  VALID_LastName,
		}
		_, err := service.CreateNewClient(ctx, in)
		wantErr := ErrPasswordNotSupplied
		if err != wantErr {
			t.Errorf("Expected err: %v Got : %v", wantErr, err)
		}
	})

	t.Run("FirstName_not_supplied", func(t *testing.T) {
		t.Cleanup(mock.Teardown)
		ctx := context.Background()
		in := RegisterUserRequest{
			Email:    VALID_EMAIL_A,
			Username: VALID_USERNAME_A,
			Password: VALID_PASSWORD,
			LastName: VALID_LastName,
		}
		_, err := service.CreateNewClient(ctx, in)
		wantErr := ErrFirstNameNotSupplied
		if err != wantErr {
			t.Errorf("Expected err: %v Got : %v", wantErr, err)
		}
	})
}

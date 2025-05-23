package auth

import (
	"context"
	"os"
	"testing"

	"b2b.nati011.github.com/internal/adapter/secondary/application/auth/provider"
<<<<<<< HEAD
	"b2b.nati011.github.com/internal/core/application/email"
=======
>>>>>>> 8bacbbe2 (- resolve weird issues)
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

<<<<<<< HEAD
var service Provider
=======
var testContainer TestContainer
>>>>>>> 8bacbbe2 (- resolve weird issues)
var mock = provider.NewMockAuthProvider()

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
<<<<<<< HEAD
	service = NewAuthService(&mock, email.NewTestContainer().EmailService)
}

func Test_CreateClient_happyPath(t *testing.T) {
	t.Cleanup(mock.Teardown)
=======
	testContainer = NewTestContainer()
}

func teardown() {
	testContainer.teardown()
}

func Test_CreateClient_happyPath(t *testing.T) {
	t.Cleanup(teardown)
>>>>>>> 8bacbbe2 (- resolve weird issues)
	ctx := context.Background()
	in := RegisterUserRequest{
		Username:  VALID_USERNAME_A,
		Password:  VALID_PASSWORD,
		FirstName: VALID_FirstName,
		LastName:  VALID_LastName,
		Email:     VALID_EMAIL_A,
	}
<<<<<<< HEAD
	_, err := service.CreateNewClientWithPassword(ctx, in)
=======
	_, err := testContainer.Service.CreateNewClientWithPassword(ctx, &in)
>>>>>>> 8bacbbe2 (- resolve weird issues)
	if err != nil {
		t.Errorf("Failed to create err: %v", err)
	}
}

func Test_CreateClient_UnhappyPath(t *testing.T) {
	t.Run("email_not_supplied", func(t *testing.T) {
<<<<<<< HEAD
		t.Cleanup(mock.Teardown)
=======
		t.Cleanup(teardown)
>>>>>>> 8bacbbe2 (- resolve weird issues)
		ctx := context.Background()
		in := RegisterUserRequest{
			Username:  VALID_USERNAME_A,
			Password:  VALID_PASSWORD,
			FirstName: VALID_FirstName,
			LastName:  VALID_LastName,
		}
<<<<<<< HEAD
		_, err := service.CreateNewClientWithPassword(ctx, in)
=======
		_, err := testContainer.Service.CreateNewClientWithPassword(ctx, &in)
>>>>>>> 8bacbbe2 (- resolve weird issues)
		wantErr := ErrEmailNotSupplied
		if err != wantErr {
			t.Errorf("Expected err: %v Got : %v", wantErr, err)
		}
	})

	t.Run("password_not_supplied", func(t *testing.T) {
<<<<<<< HEAD
		t.Cleanup(mock.Teardown)
=======
		t.Cleanup(teardown)
>>>>>>> 8bacbbe2 (- resolve weird issues)
		ctx := context.Background()
		in := RegisterUserRequest{
			Email:     VALID_EMAIL_A,
			Username:  VALID_USERNAME_A,
			FirstName: VALID_FirstName,
			LastName:  VALID_LastName,
		}
<<<<<<< HEAD
		_, err := service.CreateNewClientWithPassword(ctx, in)
=======
		_, err := testContainer.Service.CreateNewClientWithPassword(ctx, &in)
>>>>>>> 8bacbbe2 (- resolve weird issues)
		wantErr := ErrPasswordNotSupplied
		if err != wantErr {
			t.Errorf("Expected err: %v Got : %v", wantErr, err)
		}
	})

	t.Run("FirstName_not_supplied", func(t *testing.T) {
<<<<<<< HEAD
		t.Cleanup(mock.Teardown)
=======
		t.Cleanup(teardown)
>>>>>>> 8bacbbe2 (- resolve weird issues)
		ctx := context.Background()
		in := RegisterUserRequest{
			Email:    VALID_EMAIL_A,
			Username: VALID_USERNAME_A,
			Password: VALID_PASSWORD,
			LastName: VALID_LastName,
		}
<<<<<<< HEAD
		_, err := service.CreateNewClientWithPassword(ctx, in)
=======
		_, err := testContainer.Service.CreateNewClientWithPassword(ctx, &in)
>>>>>>> 8bacbbe2 (- resolve weird issues)
		wantErr := ErrFirstNameNotSupplied
		if err != wantErr {
			t.Errorf("Expected err: %v Got : %v", wantErr, err)
		}
	})
}

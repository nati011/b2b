package auth

import (
	"testing"

	authDTO "b2b.nati011.github.com/cmd/auth/model/dto"
)

type AuthTestContainer struct {
}

func (*AuthTestContainer) NewTestContainer() (*AuthTestContainer, error) {
	return &AuthTestContainer{}, nil
}

type testCaseAuth struct {
	input authDTO.RegisterUserRequest
	want  authDTO.RegisterUserResponse
}

func TestPassCreateClient(t *testing.T) {
	user := authDTO.RegisterUserRequest{

		Username:        "expired_pineapple",
		Password:        "test@123",
		ConfirmPassword: "test@123",
		FullName:        "ruth tirusew",
		Email:           "ruthtirusew944@gmail.com",
	}

	userRegistrationSuccessResponse := authDTO.RegisterUserResponse{
		Username: "expired_pineapple",
		Message:  "Ahoy!",
	}

	testCase := testCaseAuth{
		input: user,
		want:  userRegistrationSuccessResponse,
	}

	got := CreateClient(user)

	if got != testCase.want {
		t.Errorf("want:%v \n got:%v", testCase.want, got)
	}
}

func TestFailDuplicateClientUsername(t *testing.T) {
	users := []authDTO.RegisterUserRequest{
		{
			Username:        "expired_pineapple",
			Password:        "test@123",
			ConfirmPassword: "test@123",
			FullName:        "ruth tirusew",
			Email:           "ruthtirusew944@gmail.com",
		},
		{
			Username:        "expired_pineapple",
			Password:        "test@123",
			ConfirmPassword: "test@123",
			FullName:        "ruth tirusew",
			Email:           "ruthtirusew30@gmail.com",
		},
	}

	userRegistrationFailureResponse := authDTO.RegisterUserResponse{
		Username: "",
		Message:  "Oopsy, username is already taken",
	}

	// testCase :=
}

func TestFailDuplicateClientEmail(t *testing.T) {
	users := []authDTO.RegisterUserRequest{
		{
			Username:        "expired_pineapple",
			Password:        "test@123",
			ConfirmPassword: "test@123",
			FullName:        "ruth tirusew",
			Email:           "ruthtirusew944@gmail.com",
		},
		{
			Username:        "delilah",
			Password:        "test@123",
			ConfirmPassword: "test@123",
			FullName:        "ruth tirusew",
			Email:           "ruthtirusew944@gmail.com",
		},
	}

	userRegistrationFailureResponse := authDTO.RegisterUserResponse{
		Username: "",
		Message:  "oopsy, username is already taken",
	}

	// testCase :=
}

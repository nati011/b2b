package auth

import (
	"testing"

	authDTO "b2b.nati011.github.com/cmd/auth/model/dto"
)

type testCaseAuth struct {
	input authDTO.RegisterUserRequest
	want  authDTO.RegisterUserResponse
}

func TestCreateClient_happyPath(t *testing.T) {
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

	got, err := CreateClient(testCase.input)
	if err != nil {
		t.Errorf("failed to create client err: %v", err)
	}
	if got != testCase.want {
		t.Errorf("want:%v \n got:%v", testCase.want, got)
	}
}

func TestCreateClient_unhappyPath(t *testing.T) {
	t.Run("duplicateUsername", func(t *testing.T) {

		//prepare for testing duplicate username
		ua := authDTO.RegisterUserRequest{
			Username:        "expired_pineapple",
			Password:        "test@123",
			ConfirmPassword: "test@123",
			FullName:        "ruth tirusew",
			Email:           "ruthtirusew944@gmail.com",
		}

		CreateClient(ua)
		ub := authDTO.RegisterUserRequest{
			Username:        "expired_pineapple",
			Password:        "test@123",
			ConfirmPassword: "test@123",
			FullName:        "ruth tirusew",
			Email:           "ruthtirusew30@gmail.com",
		}

		want := authDTO.RegisterUserResponse{
			Username: "",
			Message:  "Oopsy, username is already taken",
		}

		got, err := CreateClient(ub)
		if err != nil {
			t.Errorf("failed to create client err: %v", err)
		}

		if got != want {
			t.Errorf("want:%v \n got:%v", want, got)
		}

	})

	t.Run("duplicateEmail", func(t *testing.T) {

		//prepare for testing duplicate username
		ua := authDTO.RegisterUserRequest{
			Username:        "expired_pineapple",
			Password:        "test@123",
			ConfirmPassword: "test@123",
			FullName:        "ruth tirusew",
			Email:           "ruthtirusew944@gmail.com",
		}

		CreateClient(ua)
		ub := authDTO.RegisterUserRequest{
			Username:        "expired_pineapple",
			Password:        "test@123",
			ConfirmPassword: "test@123",
			FullName:        "ruth tirusew",
			Email:           "ruthtirusew30@gmail.com",
		}

		want := authDTO.RegisterUserResponse{
			Username: "",
			Message:  "Oopsy, email is already taken",
		}

		got, err := CreateClient(ub)
		if err != nil {
			t.Errorf("failed to create client err: %v", err)
		}

		if got != want {
			t.Errorf("want:%v \n got:%v", want, got)
		}

	})
}

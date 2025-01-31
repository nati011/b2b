package test

import (
	"testing"

	authDTO "b2b.nati011.github.com/cmd/auth/model/dto"
)

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

	got, err := authTestContainer.authService.CreateClient(testCase.input)
	if err != nil {
		t.Errorf("failed to create client err: %v", err)
	}
	if got != testCase.want {
		t.Errorf("want:%v \n got:%v", testCase.want, got)
	}
}

func TestFailDuplicateUsernameCreateClient(t *testing.T) {
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

	want := authDTO.RegisterUserResponse{
		Username: "",
		Message:  "Oopsy, username is already taken",
	}

	//prepare for testing duplicate username
	authTestContainer.authService.CreateClient(users[0])

	got, err := authTestContainer.authService.CreateClient(users[1])
	if err != nil {
		t.Errorf("failed to create client err: %v", err)
	}

	if got != want {
		t.Errorf("want:%v \n got:%v", want, got)
	}
}

// func TestFailDuplicateClientEmail(t *testing.T) {
// 	users := []authDTO.RegisterUserRequest{
// 		{
// 			Username:        "expired_pineapple",
// 			Password:        "test@123",
// 			ConfirmPassword: "test@123",
// 			FullName:        "ruth tirusew",
// 			Email:           "ruthtirusew944@gmail.com",
// 		},
// 		{
// 			Username:        "delilah",
// 			Password:        "test@123",
// 			ConfirmPassword: "test@123",
// 			FullName:        "ruth tirusew",
// 			Email:           "ruthtirusew944@gmail.com",
// 		},
// 	}

// 	userRegistrationFailureResponse := authDTO.RegisterUserResponse{
// 		Username: "",
// 		Message:  "oopsy, username is already taken",
// 	}

// }

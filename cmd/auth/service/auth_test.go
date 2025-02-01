package auth

import (
	"testing"

	authDTO "b2b.nati011.github.com/cmd/auth/model/dto"
)

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

	in := user
	want := userRegistrationSuccessResponse

	//TODO: mock provider
	got, err := CreateClient(in)
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
		//TODO: mock provider
		CreateClient(in_a)

		in_b := authDTO.RegisterUserRequest{
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
		//TODO: mock provider
		got, err := CreateClient(in_b)
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

		//TODO: mock provider
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

		//TODO: mock provider
		got, err := CreateClient(ub)
		if err != nil {
			t.Errorf("failed to create client err: %v", err)
		}

		if got != want {
			t.Errorf("want:%v \n got:%v", want, got)
		}

	})
}

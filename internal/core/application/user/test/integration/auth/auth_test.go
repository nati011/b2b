package auth

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"b2b.nati011.github.com/internal/core/application/auth"
	"b2b.nati011.github.com/internal/core/application/user"
	db_test_container "b2b.nati011.github.com/internal/core/util/test_container/db"
)

var testContainer user.TestContainer

var db *sql.DB

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	db = db_test_container.Setup()
	testContainer = user.NewIntegrationTestContainer(
		db,
	)
}

func teardown() {
	testContainer.TeardownIntegrationTestContainer()
	db_test_container.Teardown(db)
}

func Test_create_auth_client_upon_user_registration(t *testing.T) {
	t.Cleanup(teardown)
	ctx := context.Background()
	//setup
	parsedTime, _ := time.Parse("2006-01-02 15:04:05", "2024-09-19 14:00:00")
	in := user.CreateRequest{
		FirstName: "natnael asefa",
		LastName:  "jemaneh",
		Email:     "natnaeljemaneh001@gmail.com",
		Phone:     "+251949184879",
		Username:  "test",
		DOB:       parsedTime,

		ExternalId: "123",
	}
	_, err := testContainer.UserService.Create(ctx, &in)
	if err != nil {
		t.Fatalf("Failed to create err: %v", err)
	}

	//verify
	_, err = testContainer.AuthService.ClientLogin(ctx, auth.LoginUserRequest{
		Email:    "natnaeljemaneh001@gmail.com",
		Password: "test",
	})
	if err == auth.ErrFailedToLogin {
		t.Fatalf("failed to login err %v", err)
	}
}

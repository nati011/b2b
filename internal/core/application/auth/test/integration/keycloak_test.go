package integration

import (
	"context"
	"os"
	"testing"

	provider "b2b.nati011.github.com/internal/adapter/secondary/application/auth/provider"
	"b2b.nati011.github.com/internal/core/application/auth"
	port "b2b.nati011.github.com/internal/port/application/auth/provider"

	keycloak "github.com/stillya/testcontainers-keycloak"
)

var keycloakContainer *keycloak.KeycloakContainer
var KeycloakProvider port.Provider
var authService port.Provider

const (
	VALID_PASSWORD   = "test@123"
	VALID_FIRST_NAME = "Ruth"
	VALID_LAST_NAME  = "T"
	VALID_EMAIL_A    = "ruthtirusew944@gmail.com"
	VALID_EMAIL_B    = "ruthtirusew388@gmail.com"
	VALID_USERNAME_A = "expired_pineapple"
	VALID_USERNAME_B = "delila"

	//INVALID
	INVALID_username  = ""
	INVALID_PASSWORD  = ""
	INVALID_FirstName = ""
	INVALID_EMAIL     = ""
)

func Test_Timeout(t *testing.T) {

}

func Test_CreateClient_happyPath(t *testing.T) {
	ctx := context.Background()
	user := port.RegisterUserRequest{
		Username:        VALID_USERNAME_A,
		Password:        VALID_PASSWORD,
		ConfirmPassword: VALID_PASSWORD,
		FirstName:       VALID_FIRST_NAME,
		LastName:        VALID_LAST_NAME,
		Email:           VALID_EMAIL_A,
	}

	userRegistrationSuccessResponse := port.RegisterUserResponse{
		Username: VALID_USERNAME_A,
	}

	in := user
	want := userRegistrationSuccessResponse

	got, err := authService.CreateNewClient(ctx, in)
	if err != nil {
		t.Fatalf("Failed to create client err: %v", err)
	}
	if got != want {
		t.Errorf("Expected: %v, Got: %v", want, got)
	}
	// t.Cleanup(teardown)
}

func Test_CreateClient_UnhappyPath(t *testing.T) {
	ctx := context.Background()
	t.Run("DuplicateUsername", func(t *testing.T) {
		// t.Cleanup(teardown)

		/* create a user with some username x and attempt
		to create another user with the same username */
		in_a := port.RegisterUserRequest{
			FirstName:       VALID_FIRST_NAME,
			LastName:        VALID_LAST_NAME,
			Username:        VALID_USERNAME_A,
			Password:        VALID_PASSWORD,
			ConfirmPassword: VALID_PASSWORD,
			Email:           VALID_EMAIL_A,
		}

		_, err := authService.CreateNewClient(ctx, in_a)
		if err != nil {
			t.Fatalf("Failed to create client err: %v", err)
		}

		in_b := port.RegisterUserRequest{
			Username:        VALID_USERNAME_A,
			Password:        VALID_PASSWORD,
			ConfirmPassword: VALID_PASSWORD,
			LastName:        VALID_LAST_NAME,
			Email:           VALID_EMAIL_B,
		}

		_, err = authService.CreateNewClient(ctx, in_b)
		wantErr := auth.ErrUsernameTaken
		if err != wantErr {
			t.Errorf("Expected err: %v, Got: %v", wantErr, err)
		}

	})

	t.Run("DuplicateEmail", func(t *testing.T) {
		// t.Cleanup(teardown)

		/* create a user with some email x and attempt
		to create another user with the same email */
		ua := port.RegisterUserRequest{
			Username:        VALID_USERNAME_A,
			Password:        VALID_PASSWORD,
			ConfirmPassword: VALID_PASSWORD,
			FirstName:       VALID_FIRST_NAME,
			Email:           VALID_EMAIL_A,
		}

		authService.CreateNewClient(ctx, ua)
		ub := port.RegisterUserRequest{
			FirstName:       VALID_FIRST_NAME,
			LastName:        VALID_LAST_NAME,
			Username:        VALID_USERNAME_B,
			Password:        VALID_PASSWORD,
			ConfirmPassword: VALID_PASSWORD,
			Email:           VALID_EMAIL_A,
		}

		_, err := authService.CreateNewClient(ctx, ub)
		wantErr := auth.ErrEmailTaken
		if err != wantErr {
			t.Errorf("Expected err: %v, Got: %v", wantErr, err)
		}
	})
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
	shutDown()
}

func setup() {
	// var err error
	// ctx := context.Background()
	// keycloakContainer, err = RunContainer(ctx)
	// if err != nil {
	// 	panic(err)
	// }

	// keycloakInstanceUrl, err := keycloakContainer.GetAuthServerURL(ctx)
	// if err != nil {
	// 	panic(err)
	// }

	// keycloakAdminClient, err := keycloakContainer.GetAdminClient(ctx)
	// if err != nil {
	// 	panic(err)
	// }
	KeycloakUsername := "ruthtirusew944@gmail.com"
	KeycloakPassword := "W>-553:F?XWXpmV"
	KeycloakRealm := "b2b"
	keycloakApplicationRealm := "b2b"
	keycloakClientId := "test"
	keycloakInstanceUrl := "https://euc1.auth.ac/auth"

	KeycloakProvider := provider.NewKeycloakProvider(
		keycloakInstanceUrl,
		KeycloakUsername,
		KeycloakPassword,
		KeycloakRealm,
		keycloakApplicationRealm,
		keycloakClientId,
	)

	authService = auth.NewAuthService(KeycloakProvider)
}

func shutDown() {
	ctx := context.Background()
	err := keycloakContainer.Terminate(ctx)
	if err != nil {
		panic(err)
	}
}

func teardown() {
	shutDown()
	setup()
}

const (
	KEYCLOAK_VERSION        = "keycloak/keycloak:24.0"
	KEYCLOAK_ADMIN_USERNAME = "admin"
	KEYCLOAK_ADMIN_PASSWORD = "admin"
	// KEYCLOAK_ADMIN_IMPORTFILE   = "../testdata/realm-export.json"
	KEYCLOAK_ADMIN_CONTEXT_PATH = "/auth"
)

func RunContainer(ctx context.Context) (*keycloak.KeycloakContainer, error) {
	return keycloak.Run(ctx,
		KEYCLOAK_VERSION,
		keycloak.WithContextPath(KEYCLOAK_ADMIN_CONTEXT_PATH),
		// keycloak.WithRealmImportFile(KEYCLOAK_ADMIN_IMPORTFILE),
		keycloak.WithAdminUsername(KEYCLOAK_ADMIN_USERNAME),
		keycloak.WithAdminPassword(KEYCLOAK_ADMIN_PASSWORD),
	)
}

package keycloak

import (
	"context"
	"os"
	"testing"

	"b2b.nati011.github.com/internal/adapter/secondary/application/auth/provider"
	"b2b.nati011.github.com/internal/core/application/auth"
	"b2b.nati011.github.com/internal/core/application/email"
	"b2b.nati011.github.com/internal/core/application/role"

	keycloak "github.com/stillya/testcontainers-keycloak"
)

var keycloakContainer *keycloak.KeycloakContainer
var KeycloakProvider auth.Provider
var authService auth.Provider

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
	t.Cleanup(teardown)
	ctx := context.Background()
	in := auth.RegisterUserRequest{
		Username:  VALID_USERNAME_A,
		Password:  VALID_PASSWORD,
		FirstName: VALID_FIRST_NAME,
		LastName:  VALID_LAST_NAME,
		Email:     VALID_EMAIL_A,
	}

	_, err := authService.CreateNewClientWithPassword(ctx, &in)
	if err != nil {
		t.Fatalf("Failed to create client err: %v", err)
	}
}

func Test_CreateClient_UnhappyPath(t *testing.T) {
	ctx := context.Background()
	t.Run("DuplicateUsername", func(t *testing.T) {
		t.Cleanup(teardown)

		/* create a user with some username x and attempt
		to create another user with the same username */
		in_a := auth.RegisterUserRequest{
			FirstName: VALID_FIRST_NAME,
			LastName:  VALID_LAST_NAME,
			Username:  VALID_USERNAME_A,
			Password:  VALID_PASSWORD,
			Email:     VALID_EMAIL_A,
		}

		_, err := authService.CreateNewClientWithPassword(ctx, &in_a)
		if err != nil {
			t.Fatalf("Failed to create client err: %v", err)
		}

		in_b := auth.RegisterUserRequest{
			FirstName: VALID_FIRST_NAME,
			Username:  VALID_USERNAME_A,
			Password:  VALID_PASSWORD,
			LastName:  VALID_LAST_NAME,
			Email:     VALID_EMAIL_B,
		}

		_, err = authService.CreateNewClientWithPassword(ctx, &in_b)
		wantErr := auth.ErrUsernameTaken
		if err != wantErr {
			t.Errorf("Expected err: %v, Got: %v", wantErr, err)
		}

	})

	t.Run("DuplicateEmail", func(t *testing.T) {
		t.Cleanup(teardown)

		/* create a user with some email x and attempt
		to create another user with the same email */
		ua := auth.RegisterUserRequest{
			Username:  VALID_USERNAME_A,
			Password:  VALID_PASSWORD,
			FirstName: VALID_FIRST_NAME,
			LastName:  VALID_LAST_NAME,
			Email:     VALID_EMAIL_A,
		}

		authService.CreateNewClientWithPassword(ctx, &ua)
		ub := auth.RegisterUserRequest{
			FirstName: VALID_FIRST_NAME,
			LastName:  VALID_LAST_NAME,
			Username:  VALID_USERNAME_B,
			Password:  VALID_PASSWORD,
			Email:     VALID_EMAIL_A,
		}

		_, err := authService.CreateNewClientWithPassword(ctx, &ub)
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
}

func setup() {
	var err error
	ctx := context.Background()
	keycloakContainer, err = RunContainer(ctx)
	if err != nil {
		panic(err)
	}

	keycloakInstanceUrl, err := keycloakContainer.GetAuthServerURL(ctx)
	if err != nil {
		panic(err)
	}

	keycloakAdminClient, err := keycloakContainer.GetAdminClient(ctx)
	if err != nil {
		panic(err)
	}
	// KeycloakUsername := "ruthtirusew944@gmail.com"
	// KeycloakPassword := "W>-553:F?XWXpmV"
	// KeycloakRealm := "b2b"
	// keycloakApplicationRealm := "b2b"
	// keycloakClientId := "test"
	// keycloakInstanceUrl := "https://euc1.auth.ac/auth"

	KeycloakUsername := keycloakAdminClient.Username
	KeycloakPassword := keycloakAdminClient.Password
	KeycloakRealm := keycloakAdminClient.Realm
	keycloakApplicationRealm := keycloakAdminClient.Realm
	keycloakClientId := keycloakAdminClient.ClientID

	KeycloakProvider := provider.NewKeycloakProvider(
		keycloakInstanceUrl,
		KeycloakUsername,
		KeycloakPassword,
		KeycloakRealm,
		keycloakApplicationRealm,
		keycloakClientId,
		"",
	)

	authService = auth.NewAuthService(
		KeycloakProvider,
		email.NewTestContainer().EmailService,
		role.NewTestContainer().RoleService)
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

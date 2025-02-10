package integration

import (
	"context"
	"fmt"
	"os"
	"testing"

	keycloak "github.com/stillya/testcontainers-keycloak"
)

var keycloakContainer *keycloak.KeycloakContainer

func Test_CreateNewClient(t *testing.T) {
}

func Test_ClientLogin(t *testing.T) {

}

func TestMain(m *testing.M) {
	defer func() {
		if r := recover(); r != nil {
			shutDown()
			fmt.Println("Panic")
		}
	}()
	setup()
	code := m.Run()
	shutDown()
	os.Exit(code)
}

func setup() {
	var err error
	ctx := context.Background()
	keycloakContainer, err = RunContainer(ctx)
	if err != nil {
		panic(err)
	}
}

func shutDown() {
	ctx := context.Background()
	err := keycloakContainer.Terminate(ctx)
	if err != nil {
		panic(err)
	}
}

func RunContainer(ctx context.Context) (*keycloak.KeycloakContainer, error) {
	return keycloak.Run(ctx,
		"keycloak/keycloak:24.0",
		keycloak.WithContextPath("/auth"),
		keycloak.WithRealmImportFile("../testdata/realm-export.json"),
		keycloak.WithAdminUsername("admin"),
		keycloak.WithAdminPassword("admin"),
	)
}

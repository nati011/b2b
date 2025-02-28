package distributor

import (
	"os"
	"testing"

	"b2b.nati011.github.com/internal/core/application/service/user"
	"b2b.nati011.github.com/internal/core/domain/distributor"
)

var testContainer distributor.TestContainer

var UserService user.Provider

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	testContainer = distributor.NewPackageIntegrationTestContainer()
	UserService = testContainer.UserService
}

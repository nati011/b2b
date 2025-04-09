package resource

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"b2b.nati011.github.com/internal/core/application/resource"
	test_container "b2b.nati011.github.com/internal/core/application/role/test"
	db_test_container "b2b.nati011.github.com/internal/core/util/test_container/db"
)

var container test_container.TestContainer
var db *sql.DB
var resource_id int

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func setup() {
	db = db_test_container.Setup()
	container = test_container.NewIntegrationTestContainer(db)
	ctx := context.Background()
	resource_id, _ = container.ResourceService.Create(ctx, &resource.CreateRequest{
		Action: "test",
		Name:   "test",
	})
}

func teardown() {
	container.TeardownIntegrationTestContainer(db)
	db_test_container.Teardown(db)
}

func Test_Timeout(t *testing.T) {
}

func Test_Read(t *testing.T) {
	t.Run("get_all_resources", func(t *testing.T) {
		t.Cleanup(teardown)
		setup()
	})
}

func Test_Write(t *testing.T) {
	t.Run("add_resource", func(t *testing.T) {
		t.Cleanup(teardown)
		setup()
	})

	t.Run("remove_resource", func(t *testing.T) {
		t.Cleanup(teardown)
		setup()
	})
}

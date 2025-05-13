package db

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"b2b.nati011.github.com/internal/core/application/template"
	test_template "b2b.nati011.github.com/internal/core/application/template/test"
	db_test_container "b2b.nati011.github.com/internal/core/util/test_container/db"
)

var testContainer test_template.TestContainer
var db *sql.DB

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	db = db_test_container.Setup()
	testContainer = test_template.NewIntegrationTestContainer(db)
}

func teardown() {
	testContainer.Teardown()
	db_test_container.Teardown(db)
}

func Test_Timeout(t *testing.T) {
}

func Test_write(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		id, err := testContainer.TemplateService.Create(ctx, &template.CreateRequest{
			Name:         "test",
			HtmlTemplate: "test",
		})
		if err != nil {
			t.Fatalf("failed to create template %v", err)
		}

		//check
		got, err := testContainer.TemplateService.Get(ctx, id)
		if err != nil {
			t.Fatalf("failed to get err: %v", err)
		}

		if got.Id != id {
			t.Errorf("expected err: %v Got: %v", id, got.Id)
		}
	})

	// t.Run("update", func(t *testing.T) {
	// 	t.Cleanup(teardown)
	// 	ctx := context.Background()
	// })
}

func Test_read(t *testing.T) {
	t.Run("getByName", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		id, err := testContainer.TemplateService.Create(ctx, &template.CreateRequest{
			Name:         "test",
			HtmlTemplate: "test",
		})
		if err != nil {
			t.Fatalf("failed to create template %v", err)
		}

		got, err := testContainer.TemplateService.GetByName(ctx, "test")
		if err != nil {
			t.Fatalf("failed to get err: %v", err)
		}

		if got.Name != "test" {
			t.Errorf("expected err: %v Got: %v", id, got.Id)
		}
	})

	t.Run("getById", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		id, err := testContainer.TemplateService.Create(ctx, &template.CreateRequest{
			Name:         "test",
			HtmlTemplate: "test",
		})
		if err != nil {
			t.Fatalf("failed to create template %v", err)
		}

		got, err := testContainer.TemplateService.Get(ctx, id)
		if err != nil {
			t.Fatalf("failed to get err: %v", err)
		}

		if got.Id != id {
			t.Errorf("expected err: %v Got: %v", id, got.Id)
		}
	})

	t.Run("getAll", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		_, err := testContainer.TemplateService.Create(ctx, &template.CreateRequest{
			Name:         "test",
			HtmlTemplate: "test",
		})
		if err != nil {
			t.Fatalf("failed to create template %v", err)
		}

		got, err := testContainer.TemplateService.GetAll(ctx)
		if err != nil {
			t.Fatalf("failed to get err: %v", err)
		}

		if len(got.List) != 1 {
			t.Errorf("expected len: %v Got: %v", len(got.List))
		}
	})
}

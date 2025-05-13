package template

import (
	"context"
	"os"
	"testing"
)

var container TestContainer

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	container = NewTestContainer()
}

func Test_Create_happyPath(t *testing.T) {
	t.Cleanup(container.Teardown)
	ctx := context.Background()
	id, err := container.TemplateService.Create(ctx, &CreateRequest{
		Name:         "test",
		HtmlTemplate: "test",
	})
	if err != nil {
		t.Fatalf("failed to create template %v", err)
	}

	//check
	got, err := container.TemplateService.Get(ctx, id)
	if err != nil {
		t.Fatalf("failed to get err: %v", err)
	}

	if got.Id != id {
		t.Errorf("expected err: %v Got: %v", id, got.Id)
	}
}

func Test_Create_unhappyPath(t *testing.T) {
	t.Run("emptyHtmlTemplate", func(t *testing.T) {
		t.Cleanup(container.Teardown)
		ctx := context.Background()
		_, err := container.TemplateService.Create(ctx, &CreateRequest{
			Name:         "test",
			HtmlTemplate: "",
		})
		wantErr := ErrInvalidTemplate
		if err != wantErr {
			t.Errorf("expected err: %v got err: %v", wantErr, err)
		}
	})

	t.Run("emptyName", func(t *testing.T) {
		t.Cleanup(container.Teardown)
		ctx := context.Background()
		_, err := container.TemplateService.Create(ctx, &CreateRequest{
			Name:         "",
			HtmlTemplate: "test",
		})
		wantErr := ErrInvalidName
		if err != wantErr {
			t.Errorf("expected err: %v got err: %v", wantErr, err)
		}
	})

	t.Run("duplicateName", func(t *testing.T) {
		t.Cleanup(container.Teardown)
		ctx := context.Background()
		_, err := container.TemplateService.Create(ctx, &CreateRequest{
			Name:         "test",
			HtmlTemplate: "test",
		})
		if err != nil {
			t.Fatalf("failed to create template %v", err)
		}

		_, err = container.TemplateService.Create(ctx, &CreateRequest{
			Name:         "test",
			HtmlTemplate: "test12",
		})
		wantErr := ErrDuplicateName
		if err != wantErr {
			t.Errorf("expected err: %v got err: %v", wantErr, err)
		}
	})
}

func Test_Get_happyPath(t *testing.T) {
	t.Cleanup(container.Teardown)
	ctx := context.Background()
	id, err := container.TemplateService.Create(ctx, &CreateRequest{
		Name:         "test",
		HtmlTemplate: "test",
	})
	if err != nil {
		t.Fatalf("failed to create template %v", err)
	}

	got, err := container.TemplateService.Get(ctx, id)
	if err != nil {
		t.Fatalf("failed to get err: %v", err)
	}

	if got.Id != id {
		t.Errorf("expected err: %v Got: %v", id, got.Id)
	}
}

func Test_Get_unhappyPath(t *testing.T) {
	t.Run("idNotFound", func(t *testing.T) {
		t.Cleanup(container.Teardown)
		ctx := context.Background()
		_, err := container.TemplateService.Get(ctx, 99)
		wantErr := ErrIdNotFound
		if err != wantErr {
			t.Errorf("expected err: %v got err: %v", wantErr, err)
		}
	})

	// t.Run("nameNotFound", func(t *testing.T) {
	// 	t.Cleanup(container.Teardown)
	// 	ctx := context.Background()
	// 	_, err := container.TemplateService.getByName(ctx, "non_existing")
	// 	wantErr := ErrNameNotFound
	// 	if err != wantErr {
	// 		t.Errorf("expected err: %v got err: %v", wantErr, err)
	// 	}
	// })
}

func Test_Get_All_happyPath(t *testing.T) {
	t.Cleanup(container.Teardown)
	ctx := context.Background()
	_, err := container.TemplateService.Create(ctx, &CreateRequest{
		Name:         "test",
		HtmlTemplate: "test",
	})
	if err != nil {
		t.Fatalf("failed to create template %v", err)
	}

	got, err := container.TemplateService.GetAll(ctx)
	if err != nil {
		t.Fatalf("failed to get err: %v", err)
	}
	wantLen := 1
	if len(got.List) != wantLen {
		t.Errorf("Expected len: %v Got: %v", wantLen, len(got.List))
	}
}

func Test_Get_All_unhappyPath(t *testing.T) {
	t.Run("notFound", func(t *testing.T) {
		t.Cleanup(container.Teardown)
		ctx := context.Background()
		_, err := container.TemplateService.GetAll(ctx)
		wantErr := ErrEmptyGetContent
		if err != wantErr {
			t.Errorf("Expected err: %v Got: %v", wantErr, err)
		}
	})
}

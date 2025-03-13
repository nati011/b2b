package template

import (
	"context"
	"os"
	"testing"

	db "b2b.nati011.github.com/internal/adapter/secondary/application/email-template/db"
)

var templateService Provider

func Test_Create_happyPath(t *testing.T) {
	ctx := context.Background()
	in := CreateRequest{
		Name:         "simple",
		HtmlTemplate: "test",
	}
	want := CreateResponse{
		Name:    "simple",
		Message: SUCCESS_MESSAGE,
	}

	got, err := templateService.Create(ctx, in)
	if err != nil {
		t.Errorf("Failed to create template err: %q", err)
	}
	if got != want {
		t.Errorf("Expected: %q Got: %q", want, got)
	}
}

func Test_Create_unhappyPath(t *testing.T) {
	t.Run("emptyHtmlTemplate", func(t *testing.T) {
		ctx := context.Background()
		in := CreateRequest{
			Name:         "simple",
			HtmlTemplate: "",
		}
		want := CreateResponse{
			Name:    "",
			Message: ErrSysInvalidTemplate.Error(),
		}

		got, err := templateService.Create(ctx, in)
		if err != ErrSysInvalidTemplate {
			t.Errorf("Failed to create template err: %q", err)
		}
		if got != want {
			t.Errorf("Expected: %q Got: %q", want, got)
		}
	})

	t.Run("emptyName", func(t *testing.T) {
		ctx := context.Background()
		in := CreateRequest{
			Name:         "",
			HtmlTemplate: "test",
		}
		want := CreateResponse{
			Name:    "",
			Message: ErrSysInvalidName.Error(),
		}

		got, err := templateService.Create(ctx, in)
		if err != ErrSysInvalidName {
			t.Errorf("Failed to create template err: %v", err)
		}
		if got != want {
			t.Errorf("Expected: %q Got: %q", want, got)
		}
	})

	t.Run("duplicateName", func(t *testing.T) {
		//init
		ctx := context.Background()
		_, err := templateService.Create(ctx, CreateRequest{
			Name:         "simple",
			HtmlTemplate: "test",
		})
		if err != nil {
			t.Errorf("Failed to create template err: %v", err)
		}

		in := CreateRequest{
			Name:         "simple",
			HtmlTemplate: "test",
		}
		want := CreateResponse{
			Name:    "simple",
			Message: ErrSysDuplicateName.Error(),
		}

		got, err := templateService.Create(ctx, in)
		if err != ErrSysDuplicateName {
			t.Errorf("Failed to create template err: %v", err)
		}
		if got != want {
			t.Errorf("Expected: %q Got: %q", want, got)
		}
	})
}

func Test_Update_happyPath(t *testing.T) {
}

func Test_Update_unhappyPath(t *testing.T) {
}

func Test_Get_happyPath(t *testing.T) {
	//init
	ctx := context.Background()
	_, err := templateService.Create(ctx, CreateRequest{
		Name:         "simple",
		HtmlTemplate: "test",
	})
	if err != nil {
		t.Errorf("Failed to create template err: %v", err)
	}
	in := "simple"
	want := GetResponse{
		Name:         "simple",
		HtmlTemplate: "test",
	}
	got, err := templateService.Get(ctx, in)
	if err != nil {
		t.Errorf("Failed to fetch email err: %v", err)
	}
	if got != want {
		t.Errorf("Expected: %v Want: %v", want, got)
	}
}

func Test_Get_unhappyPath(t *testing.T) {
	t.Run("notFound", func(t *testing.T) {
		ctx := context.Background()
		in := "simple"
		want := GetResponse{}
		got, err := templateService.Get(ctx, in)
		if err != nil {
			t.Errorf("Failed to fetch email err: %v", err)
		}
		if got != want {
			t.Errorf("Expected: %v Want: %v", want, got)
		}
	})
}

func Test_Get_All_happyPath(t *testing.T) {
	// init
	ctx := context.Background()
	_, err := templateService.Create(ctx, CreateRequest{
		Name:         "simple",
		HtmlTemplate: "test",
	})
	if err != nil {
		t.Errorf("Failed to create template err: %v", err)
	}

	want := GetAllResponse{
		List: []GetResponse{
			{
				Name:         "simple",
				HtmlTemplate: "test",
			},
		},
	}
	got, err := templateService.GetAll(ctx)
	if err != nil {
		t.Errorf("Failed to fetch email err: %v", err)
	}
	for _, i := range got.List {
		if i.Name != want.List[0].Name || i.HtmlTemplate != want.List[0].HtmlTemplate {
			t.Errorf("Expected: %v Want: %v", want, got)
		}
	}
}

func Test_Get_All_unhappyPath(t *testing.T) {
	t.Run("notFound", func(t *testing.T) {
		ctx := context.Background()
		want := GetAllResponse{}
		got, err := templateService.GetAll(ctx)
		if err != nil {
			t.Errorf("Failed to fetch email err: %v", err)
		}
		if len(got.List) != 0 {
			t.Errorf("Expected: %v Want: %v", want, got)
		}
	})
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	templateService = NewTemplateService(db.NewMock())
}

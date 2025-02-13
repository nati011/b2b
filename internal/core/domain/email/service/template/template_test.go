package email

import (
	"os"
	"testing"

	provider "b2b.nati011.github.com/internal/core/domain/email/provider/db/template"
)

var templateService Templer

func Test_Create_happyPath(t *testing.T) {
	in := CreateRequest{
		Name:         "simple",
		HtmlTemplate: "test",
	}
	want := CreateResponse{
		Name:    "simple",
		Message: SUCCESS_MESSAGE,
	}

	got, err := templateService.Create(in)
	if err != nil {
		t.Errorf("Failed to create template err: %q", err)
	}
	if got != want {
		t.Errorf("Expected: %q Got: %q", want, got)
	}
}

func Test_Create_unhappyPath(t *testing.T) {
	t.Run("emptyHtmlTemplate", func(t *testing.T) {
		in := CreateRequest{
			Name:         "simple",
			HtmlTemplate: "",
		}
		want := CreateResponse{
			Name:    "",
			Message: ErrSysInvalidTemplate.Error(),
		}

		got, err := templateService.Create(in)
		if err != ErrSysInvalidTemplate {
			t.Errorf("Failed to create template err: %q", err)
		}
		if got != want {
			t.Errorf("Expected: %q Got: %q", want, got)
		}
	})

	t.Run("emptyName", func(t *testing.T) {
		in := CreateRequest{
			Name:         "",
			HtmlTemplate: "test",
		}
		want := CreateResponse{
			Name:    "",
			Message: ErrSysInvalidName.Error(),
		}

		got, err := templateService.Create(in)
		if err != ErrSysInvalidName {
			t.Errorf("Failed to create template err: %v", err)
		}
		if got != want {
			t.Errorf("Expected: %q Got: %q", want, got)
		}
	})

	t.Run("duplicateName", func(t *testing.T) {
		//init
		_, err := templateService.Create(CreateRequest{
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

		got, err := templateService.Create(in)
		if err != ErrSysDuplicateName {
			t.Errorf("Failed to create template err: %q", err)
		}
		if got != want {
			t.Errorf("Expected: %q Got: %q", want, got)
		}
	})
}

func Test_Get_happyPath(t *testing.T) {
	//init
	_, err := templateService.Create(CreateRequest{
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
	got, err := templateService.Get(in)
	if err != nil {
		t.Errorf("Failed to fetch email err: %v", err)
	}
	if got != want {
		t.Errorf("Expected: %v Want: %v", want, got)
	}
}

func Test_Get_unhappyPath(t *testing.T) {
	t.Run("notFound", func(t *testing.T) {
		in := "simple"
		want := GetResponse{}
		got, err := templateService.Get(in)
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
	_, err := templateService.Create(CreateRequest{
		Name:         "simple",
		HtmlTemplate: "test",
	})
	if err != nil {
		t.Errorf("Failed to create template err: %v", err)
	}
	_, err = templateService.Create(CreateRequest{
		Name:         "simple_2",
		HtmlTemplate: "test_2",
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
	got, err := templateService.GetAll()
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
		want := GetAllResponse{}
		got, err := templateService.GetAll()
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
	templateService = NewTemplateService(&provider.MockDB{})
}

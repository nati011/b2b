package email

import (
	"os"
	"testing"
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
			HtmlTemplate: "test",
		}
		want := CreateResponse{
			Name:    "simple",
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
			Name:    "simple",
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
		want := GetTemplateResponse{}
		got, err := templateService.GetTemplate(in)
		if err != nil {
			t.Errorf("Failed to fetch email err: %v", err)
		}
		if got != want {
			t.Errorf("Expected: %v Want: %v", want, got)
		}
	})
}

func Test_Get_All_happyPath(t *testing.T) {

}

func Test_Get_All_unhappyPath(t *testing.T) {
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	templateService = NewTemplateService()
}

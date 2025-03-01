package email

import (
	"os"
	"testing"

	template "b2b.nati011.github.com/internal/core/application/email/template"
)

var service Renderer

func Test_Render_happyPath(t *testing.T) {
	in := Request{
		Name: "test",
		Args: map[string]string{
			"a": "test",
			"b": "test",
		},
	}
	want := Response{
		Name: "test",
		Text: "mock",
	}

	got, err := service.Create(&in)
	if err != nil {
		t.Errorf("Failed to create template err: %v", err)
	}
	if got != want {
		t.Errorf("Expected: %v, Got: %v", want, got)
	}
}

func Test_Render_unhappyPath(t *testing.T) {
	t.Run("templateNotFound", func(t *testing.T) {
		in := Request{
			Name: "non_existing_template",
			Args: map[string]string{
				"a": "test",
				"b": "test",
			},
		}
		want := Response{
			Name: "",
			Text: "",
		}

		got, err := service.Create(&in)
		if err != ErrSysTemplateNotFound {
			t.Errorf("Failed to render templaete err: %q", err)
		}
		if got != want {
			t.Errorf("Expected: %q Got: %q", want, got)
		}

	})
}

func TestMain(m *testing.M) {
	setup()
	c := m.Run()
	os.Exit(c)
}

func setup() {
	service = NewRenderService(template.NewMock())
}

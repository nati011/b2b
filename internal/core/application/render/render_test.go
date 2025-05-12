package render

import (
	"context"
	"os"
	"testing"

	"b2b.nati011.github.com/internal/core/application/template"
)

var testContainer TestContainer
var temp_id int

func TestMain(m *testing.M) {
	setup()
	c := m.Run()
	os.Exit(c)
}

func setup() {
	testContainer = NewTestContainer()
	ctx := context.Background()
	var err error
	temp_id, err = testContainer.TemplateService.Create(ctx, &template.CreateRequest{
		Name:         "test",
		HtmlTemplate: "{{.a}}",
	})
	if err != nil {
		panic("failed to create template")
	}
}

func Test_Render_happyPath(t *testing.T) {
	in := Request{
		TemplateId: 1,
		Args: map[string]string{
			"a": "mock"},
	}
	want := Response{
		Name: "test",
		Text: "mock",
	}

	got, err := testContainer.RenderService.Create(&in)
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
			TemplateId: 99,
			Args: map[string]string{
				"a": "test",
				"b": "test",
			},
		}
		want := Response{
			Name: "",
			Text: "",
		}

		got, err := testContainer.RenderService.Create(&in)
		if err != ErrSysTemplateNotFound {
			t.Errorf("Failed to render templaete err: %q", err)
		}
		if got != want {
			t.Errorf("Expected: %q Got: %q", want, got)
		}

	})
}

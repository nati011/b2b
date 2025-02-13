package email

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	setup()
	c := m.Run()
	os.Exit(c)
}

func setup() {
	service = NewRenderService()
}

var service *RenderService

func Test_Render_happyPath(t *testing.T) {
	in := Request{
		Name: "test",
		Args: map[string]string{
			"a": "test",
			"b": "test",
		},
	}
	want := RenderResponse{
		Name:    "test",
		Subject: "test",
		Text:    "test",
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
}

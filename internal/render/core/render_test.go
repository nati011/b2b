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
		Name: "test",
		HtmlTemplate: `
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .button {
            display: inline-block;
            background-color: #007BFF;
            color: white !important;
            padding: 10px 20px;
            text-decoration: none;
            border-radius: 5px;
        }
        .footer { margin-top: 20px; font-size: 12px; color: #666; }
    </style>
</head>
<body>
    <div class="container">
        <h2>Password Reset Request</h2>
        <p>Hello {{.username}},</p>
        <p>We received a request to reset your password. Click the button below to proceed:</p>
        <p>
            <a href="{{.reset_link}}" class="button">Reset Password</a>
        </p>
        <p>If you didn't request this, please ignore this email.</p>
        <div class="footer">
            <p>Best regards,<br>The Team</p>
        </div>
    </div>
</body>
</html>
        `,
	})
	if err != nil {
		panic("failed to create template")
	}
}

func Test_Render_happyPath(t *testing.T) {
	in := Request{
		TemplateId: 1,
		Args: map[string]interface{}{
			"username":   "John Doe",
			"reset_link": "https://example.com/reset?token=abc123"},
	}
	want := Response{
		Name: "test",
		Text: `
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .button {
            display: inline-block;
            background-color: #007BFF;
            color: white !important;
            padding: 10px 20px;
            text-decoration: none;
            border-radius: 5px;
        }
        .footer { margin-top: 20px; font-size: 12px; color: #666; }
    </style>
</head>
<body>
    <div class="container">
        <h2>Password Reset Request</h2>
        <p>Hello John Doe,</p>
        <p>We received a request to reset your password. Click the button below to proceed:</p>
        <p>
            <a href="https://example.com/reset?token=abc123" class="button">Reset Password</a>
        </p>
        <p>If you didn't request this, please ignore this email.</p>
        <div class="footer">
            <p>Best regards,<br>The Team</p>
        </div>
    </div>
</body>
</html>
        `,
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
			Args: map[string]interface{}{
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

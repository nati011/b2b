package email

import (
	"context"
	"os"
	"testing"

	"b2b.nati011.github.com/internal/core/application/template"
)

const (
	VALID_EMAIL_RECEPIENT = "ruthtirusew944@mailpit.com"
	INVALID_EMAIL_ADDR    = "ruthtirusew944"
	VALID_SUBJECT         = "test"
	VALID_TEXT            = "test"
	INVALID_TEXT          = ""
)

var testContainer TestContainer
var temp_id int

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	testContainer = *NewTestContainer()
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

func Test_SendEmail_happyPath(t *testing.T) {
	in := SendRequest{
		To:      VALID_EMAIL_RECEPIENT,
		Subject: VALID_SUBJECT,
		Args: map[string]string{
			"username":   "John Doe",
			"reset_link": "https://example.com/reset?token=abc123"},
		TemplateId: temp_id,
	}

	err := testContainer.EmailService.Send(&in)
	if err != nil {
		t.Fatalf("Failed to send err: %v", err)
	}
}

func Test_SendEmail_unhappyPath(t *testing.T) {
	t.Run("invalidRecepientEmail", func(t *testing.T) {
		in := SendRequest{
			To:      INVALID_EMAIL_ADDR,
			Subject: VALID_SUBJECT,
			Args: map[string]string{
				"username":   "John Doe",
				"reset_link": "https://example.com/reset?token=abc123"},
			TemplateId: temp_id,
		}

		err := testContainer.EmailService.Send(&in)
		wantErr := ErrReceiverAddressNotValid
		if err != wantErr {
			t.Errorf("Expected err: %v Got: %v", wantErr, err)
		}
	})
}

package integration

import (
	"context"
	"os"
	"testing"

	emailProvider "b2b.nati011.github.com/internal/adapter/secondary/application/email/smtp"
	email "b2b.nati011.github.com/internal/core/application/email"
	"b2b.nati011.github.com/internal/core/application/email/test"
	"b2b.nati011.github.com/internal/core/application/template"
	"github.com/testcontainers/testcontainers-go/modules/inbucket"
)

var inbucketContainer *inbucket.InbucketContainer
var testContainer test.TestContainer
var temp_id int

const (
	VALID_EMAIL_RECEPIENT = "ruthtirusew944@mailpit.com"
	INVALID_EMAIL_ADDR    = "ruthtirusew944"
	VALID_SUBJECT         = "test"
	VALID_TEXT            = "test"
	INVALID_TEXT          = ""
)

const (
	INBUCKET_VERSION = "inbucket/inbucket:sha-2d409bb"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	var err error
	ctx := context.Background()
	inbucketContainer, err = RunContainer(ctx)
	if err != nil {
		panic(err)
	}

	smtpPort, err := inbucketContainer.SmtpConnection(ctx)
	if err != nil {
		panic(err)
	}

	testContainer = *test.NewIntegrationTestContainer(emailProvider.NewInbucket(
		"test@gmail.com",
		smtpPort,
	))

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

func RunContainer(ctx context.Context) (*inbucket.InbucketContainer, error) {
	return inbucket.Run(ctx, INBUCKET_VERSION)
}

func Teardown() {
	ctx := context.Background()
	err := inbucketContainer.Terminate(ctx)
	if err != nil {
		panic(err)
	}
}

func Test_Send_Email(t *testing.T) {
	in := email.SendRequest{
		To:      VALID_EMAIL_RECEPIENT,
		Subject: VALID_SUBJECT,
		Args: map[string]interface{}{
			"username":   "John Doe",
			"reset_link": "https://example.com/reset?token=abc123"},
		TemplateId: temp_id,
	}
	err := testContainer.EmailService.Send(&in)
	if err != nil {
		t.Fatalf("Failed to send email err: %q", err)
	}
}

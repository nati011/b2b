package auth

import (
	"context"

	"b2b.nati011.github.com/internal/adapter/secondary/application/auth/provider"
	"b2b.nati011.github.com/internal/core/application/email"
	"b2b.nati011.github.com/internal/core/application/template"
)

func NewIntegrationAuthContainer() Provider {
	var mock = provider.NewMockAuthProvider()
	emailTestContainer := email.NewTestContainer()
	ctx := context.Background()
	_, err := emailTestContainer.TemplateService.Create(ctx, &template.CreateRequest{
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
	return NewAuthService(&mock, emailTestContainer.EmailService)
}

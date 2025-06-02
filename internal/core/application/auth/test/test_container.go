package test

import (
	"context"

	"b2b.nati011.github.com/internal/adapter/secondary/application/auth/provider"
	db_resource_provider "b2b.nati011.github.com/internal/adapter/secondary/application/resource/db"
	db_role_mock "b2b.nati011.github.com/internal/adapter/secondary/application/role/db"
	"b2b.nati011.github.com/internal/core/application/auth"
	"b2b.nati011.github.com/internal/core/application/email"
	"b2b.nati011.github.com/internal/core/application/resource"
	"b2b.nati011.github.com/internal/core/application/role"
	"b2b.nati011.github.com/internal/core/application/template"
)

type TestContainer struct {
	RoleService     role.Provider
	ResourceService resource.Provider
	Service         auth.Provider
}

func NewIntegrationTestContainer() TestContainer {
	container := TestContainer{}
	container.RoleService = role.NewRole(db_role_mock.NewMock(), container.ResourceService)
	container.ResourceService = resource.NewResource(db_resource_provider.NewMock())

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
		panic("failed to create email template")
	}

	container.Service = auth.NewAuthService(
		provider.NewMockAuthProvider(),
		emailTestContainer.EmailService,
		role.NewTestContainer().RoleService)

	return container
}

func (t *TestContainer) teardown() {

}

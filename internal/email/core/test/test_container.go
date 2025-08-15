package test

import (
	smtp "b2b.nati011.github.com/internal/adapter/secondary/application/email/smtp"
	"b2b.nati011.github.com/internal/core/application/email"
	render "b2b.nati011.github.com/internal/core/application/render"
	"b2b.nati011.github.com/internal/core/application/template"
	port "b2b.nati011.github.com/internal/port/application/email"
)

type TestContainer struct {
	RendererService render.Provider
	EmailService    email.Provider
	TemplateService template.Provider
}

func NewIntegrationTestContainer(ep port.Provider) *TestContainer {
	container := &TestContainer{}
	container.TemplateService = template.NewTestContainer().TemplateService
	container.RendererService = render.NewRenderService(container.TemplateService)
	container.EmailService = email.NewEmailService(smtp.NewMock(), container.RendererService)
	return container
}

func (t *TestContainer) Teardown() {
	t.EmailService = email.NewEmailService(smtp.NewMock(), t.RendererService)
}

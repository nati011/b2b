package email

import (
	smtp "b2b.nati011.github.com/internal/adapter/secondary/application/email/smtp"
	render "b2b.nati011.github.com/internal/core/application/render"
	"b2b.nati011.github.com/internal/core/application/template"
)

type TestContainer struct {
	RendererService render.Provider
	EmailService    Provider
	TemplateService template.Provider
}

func NewTestContainer() *TestContainer {
	container := &TestContainer{}
	container.TemplateService = template.NewTestContainer().TemplateService
	container.RendererService = render.NewRenderService(container.TemplateService)
	container.EmailService = NewEmailService(smtp.NewMock(), container.RendererService)
	return container
}

func (t *TestContainer) Teardown() {
	t.EmailService = NewEmailService(smtp.NewMock(), t.RendererService)
}

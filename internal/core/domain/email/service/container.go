package email

import (
	templateProvider "b2b.nati011.github.com/internal/core/domain/email/provider/db/template"
	smtp "b2b.nati011.github.com/internal/core/domain/email/provider/smtp"
	render "b2b.nati011.github.com/internal/core/domain/email/service/render"
	template "b2b.nati011.github.com/internal/core/domain/email/service/template"
)

type TestContainer struct {
	renderer render.Renderer
	Emailer  Emailer
}

func NewTestContainer(s smtp.Provider) *TestContainer {
	container := &TestContainer{}
	container.renderer = initTestRenderService(template.NewTemplateService(&templateProvider.MockDB{}))
	container.Emailer = initTestEmailService(s, container.renderer)
	return container
}

func initTestEmailService(s smtp.Provider, r render.Renderer) Emailer {
	return &EmailService{
		smtp:     s,
		renderer: r,
	}
}

func initTestRenderService(t template.Templer) render.Renderer {
	return render.NewMock()
}

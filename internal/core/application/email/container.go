package email

import (
	smtp "b2b.nati011.github.com/internal/adapter/secondary/application/email/smtp"
	render "b2b.nati011.github.com/internal/core/application/email/render"
)

type TestContainer struct {
	renderer render.Renderer
	Emailer  Emailer
}

func NewTestContainer(s smtp.Provider) *TestContainer {
	container := &TestContainer{}
	container.renderer = initTestRenderService()
	container.Emailer = initTestEmailService(s, container.renderer)
	return container
}

func initTestEmailService(s smtp.Provider, r render.Renderer) Emailer {
	return &EmailService{
		smtp:     s,
		renderer: r,
	}
}

func initTestRenderService() render.Renderer {
	return render.NewMock()
}

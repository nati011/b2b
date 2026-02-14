package render

import template "b2b.nati011.github.com/internal/core/application/template"

type TestContainer struct {
	RenderService   Provider
	TemplateService template.Provider
}

func NewTestContainer() TestContainer {
	c := TestContainer{}
	c.TemplateService = template.NewTestContainer().TemplateService
	c.RenderService = NewRenderService(c.TemplateService)
	return c
}

func (t *TestContainer) Teardown() {
	t.RenderService = NewRenderService(t.TemplateService)
}

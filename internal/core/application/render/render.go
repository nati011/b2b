package render

import (
	"bytes"
	"context"
	"errors"
	templ "html/template"

	template "b2b.nati011.github.com/internal/core/application/template"
)

var (
	ErrSysTemplateNotFound = errors.New("template not found")
	ErrSysUnknown          = errors.New("unknown")
)

type Request struct {
	TemplateId int
	Args       map[string]interface{}
}

type Response struct {
	Name string
	Text string
}

type Provider interface {
	Create(*Request) (Response, error)
}

type RenderService struct {
	TemplateService template.Provider
}

func NewRenderService(tp template.Provider) Provider {
	return &RenderService{
		TemplateService: tp,
	}
}

func (s RenderService) Create(r *Request) (Response, error) {
	ctx := context.Background()
	queryResp, err := s.TemplateService.Get(ctx, r.TemplateId)
	if err != nil {
		switch err {
		case template.ErrIdNotFound:
			return Response{}, ErrSysTemplateNotFound
		default:
			return Response{}, ErrSysUnknown
		}
	}

	var buf bytes.Buffer
	t, err := templ.New(queryResp.Name).Parse(queryResp.HtmlTemplate)
	if err != nil {
		return Response{}, ErrSysUnknown
	}
	//render
	err = t.Execute(&buf, r.Args)
	if err != nil {
		return Response{}, ErrSysUnknown
	}

	return Response{
		Name: queryResp.Name,
		Text: buf.String(),
	}, nil
}

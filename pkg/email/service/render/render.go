package email

import (
	"bytes"
	"context"
	"errors"
	templ "html/template"

	template "b2b.nati011.github.com/pkg/email/service/template"
)

var (
	ErrSysTemplateNotFound = errors.New("template not found")
	ErrSysUnknown          = errors.New("unknown")
)

type Request struct {
	Name string
	Args map[string]string
}

type Response struct {
	Name string
	Text string
}

type Renderer interface {
	Create(*Request) (Response, error)
}

type RenderService struct {
	TemplateService template.Templer
}

func NewRenderService(tp template.Templer) Renderer {
	return &RenderService{
		TemplateService: tp,
	}
}

func (s RenderService) Create(r *Request) (Response, error) {
	ctx := context.Background()
	queryResp, err := s.TemplateService.Get(ctx, r.Name)
	if err != nil {
		// switch err {
		// case template.ErrSysUnknown:
		// 	return Response{}, ErrSysUnknown
		// }
		return Response{}, ErrSysUnknown
	}
	emptyResp := template.GetResponse{}
	if queryResp == emptyResp {
		return Response{}, ErrSysTemplateNotFound
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

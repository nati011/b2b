package email

import (
	"bytes"
	"errors"
	templ "html/template"

	template "b2b.nati011.github.com/internal/core/domain/email/service/template"
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
	templatService template.Templer
}

func NewRenderService(tp template.Templer) Renderer {
	return &RenderService{
		templatService: tp,
	}
}

func (s RenderService) Create(r *Request) (Response, error) {
	queryResp, err := s.templatService.Get(r.Name)
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

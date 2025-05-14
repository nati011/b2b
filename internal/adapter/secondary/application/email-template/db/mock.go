package email

import (
	"context"

	port "b2b.nati011.github.com/internal/port/application/email-template"
	port_commons "b2b.nati011.github.com/internal/port/commons/db"
)

type template struct {
	Id           int
	Name         string
	HtmlTemplate string
}

type Mock struct {
	templates []template
}

func NewMock() port.DB {
	return &Mock{
		templates: []template{},
	}
}

func (m *Mock) Create(ctx context.Context, r *port.CreateRequest) (int, error) {
	newId := len(m.templates) + 1
	m.templates = append(m.templates, template{
		Id:           newId,
		Name:         r.Name,
		HtmlTemplate: r.HtmlTemplate,
	})

	return newId, nil
}

func (m *Mock) Update(ctx context.Context, r *port.UpdateRequest) error {
	updatedTemplates := []template{}
	for _, t := range m.templates {
		if t.Id == r.Id {
			updatedTemplates = append(updatedTemplates, template{
				Id:           r.Id,
				Name:         r.Name,
				HtmlTemplate: r.HtmlTemplate,
			})
		} else {
			updatedTemplates = append(updatedTemplates, template(t))
		}
	}
	m.templates = updatedTemplates
	return nil
}

func (m Mock) Get(ctx context.Context, id int) (port.GetResponse, error) {
	for _, t := range m.templates {
		if t.Id == id {
			return port.GetResponse(t), nil
		}
	}
	return port.GetResponse{}, port_commons.ErrSysNoRows
}

func (m *Mock) GetByName(ctx context.Context, name string) (port.GetResponse, error) {
	for _, t := range m.templates {
		if t.Name == name {
			return port.GetResponse(t), nil
		}
	}
	return port.GetResponse{}, port_commons.ErrSysNoRows
}

func (m *Mock) GetAll(ctx context.Context) (port.GetAllResponse, error) {
	response := port.GetAllResponse{}
	for _, i := range m.templates {
		response.List = append(response.List, port.GetResponse(i))
	}
	if len(response.List) == 0 {
		return port.GetAllResponse{}, port_commons.ErrSysNoRows
	}
	return response, nil
}

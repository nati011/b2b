package email

import "context"

type MockReaderWriter struct {
	templates []template
}

func NewMock() ReaderWriter {
	return &MockReaderWriter{
		templates: []template{},
	}
}

func (m *MockReaderWriter) Create(ctx context.Context, r *CreateRequest) (CreateResponse, error) {
	for _, i := range m.templates {
		if i.name == r.Name {
			return CreateResponse{
				Name: i.name,
			}, ErrSysDuplicateName_L1
		}
	}
	m.templates = append(m.templates, template{
		name: r.Name,
		html: r.HtmlTemplate,
	})

	return CreateResponse{
		Name: r.Name,
	}, nil
}

func (m *MockReaderWriter) Update(ctx context.Context, r *UpdateRequest) (UpdateResponse, error) {
	for _, i := range m.templates {
		if i.name == r.Name {
			i.html = r.HtmlTemplate
			return UpdateResponse{
				Name: i.name,
			}, nil
		}
	}
	return UpdateResponse{}, nil
}

func (m MockReaderWriter) Get(ctx context.Context, r string) (GetResponse, error) {
	for _, i := range m.templates {
		if i.name == r {
			return GetResponse{
				Name:         i.name,
				HtmlTemplate: i.html,
			}, nil
		}
	}
	return GetResponse{}, nil
}

func (m *MockReaderWriter) GetAll(ctx context.Context) (GetAllResponse, error) {
	return GetAllResponse{}, nil
}

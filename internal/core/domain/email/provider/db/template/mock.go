package email

type MockDB struct {
	templates []template
}

func (m *MockDB) Create(r *CreateRequest) (CreateResponse, error) {
	for _, i := range m.templates {
		if i.name == r.Name {
			return CreateResponse{}, ErrSysDuplicateName_L1
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

func (m *MockDB) Get(r string) (GetResponse, error) {
	for _, i := range m.templates {
		if i.name == r {
			return GetResponse{}, nil
		}
	}
	return GetResponse{}, nil
}

func (m *MockDB) GetAll() (GetAllResponse, error) {
	return GetAllResponse{}, nil
}

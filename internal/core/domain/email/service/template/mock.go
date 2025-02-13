package email

type Mock struct {
	templates []mockTemplates
}

type mockTemplates struct {
	Name string
	Html string
}

func NewMock() Templer {
	return &Mock{
		templates: []mockTemplates{
			{
				Name: "test",
				Html: "mock",
			},
		},
	}
}

func (t *Mock) Create(req CreateRequest) (CreateResponse, error) {
	for _, i := range t.templates {
		if i.Name == req.Name {
			return CreateResponse{
				Name:    req.Name,
				Message: SUCCESS_MESSAGE,
			}, nil
		}
	}
	return CreateResponse{}, nil
}

func (t Mock) Get(r string) (GetResponse, error) {
	for _, i := range t.templates {
		if i.Name == r {
			return GetResponse{
				Name:         r,
				HtmlTemplate: "mock",
			}, nil
		}
	}
	return GetResponse{}, nil
}

func (t Mock) GetAll() (GetAllResponse, error) {

	result := []GetResponse{}
	return GetAllResponse{
		List: result,
	}, nil
}

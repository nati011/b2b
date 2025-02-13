package email

type CreateRequest struct {
	Name         string
	HtmlTemplate string
}

type CreateResponse struct {
	Id   string
	Name string
}

type GetResponse struct {
	Name string
}

type GetTemplateResponse struct {
	Name string
}

type Templer interface {
	Create(CreateRequest) (CreateResponse, error)
	Get() (GetResponse, error)
	GetTemplate(string) (GetTemplateResponse, error)
}

type TemplateService struct {
}

func (t *TemplateService) Create(CreateRequest) (CreateResponse, error) {
	return CreateResponse{}, nil
}

func (t TemplateService) Get() (GetResponse, error) {
	return GetResponse{}, nil
}

func (t TemplateService) GetTemplate(s string) (GetTemplateResponse, error) {
	return GetTemplateResponse{}, nil
}

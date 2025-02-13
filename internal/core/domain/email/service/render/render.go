package email

type Request struct {
	Name string
	Args map[string]string
}

type RenderResponse struct {
	Name    string
	Subject string
	Text    string
}

type Renderer interface {
	Create(*Request) (RenderResponse, error)
}

type RenderService struct {
}

func NewRenderService() *RenderService {
	return &RenderService{}
}

func (s RenderService) Create(r *Request) (RenderResponse, error) {
	return RenderResponse{}, nil
}

func (s RenderService) GetAll(r *Request) (RenderResponse, error) {
	return RenderResponse{}, nil
}

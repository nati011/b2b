package render

type Mock struct {
}

func NewMock() Renderer {
	return &Mock{}
}

func (s Mock) Create(r *Request) (Response, error) {
	return Response{
		Name: "Mock",
		Text: "Mock",
	}, nil
}

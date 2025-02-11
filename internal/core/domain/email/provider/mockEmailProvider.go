package email

type MockEmailProvider struct {
}

func (m *MockEmailProvider) Send(e Request) (Response, error) {
	return Response{
		Addr: e.Addr,
	}, nil
}

func NewMockEmailProvider() *MockEmailProvider {
	return &MockEmailProvider{}
}

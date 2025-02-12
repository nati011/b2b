package email

type MockEmailProvider struct {
}

func (m *MockEmailProvider) Send(r Request) error {
	return nil
}

func NewMockEmailProvider() *MockEmailProvider {
	return &MockEmailProvider{}
}

package email

type MockEmailProvider struct {
}

func (m *MockEmailProvider) Send(e Request) error {
	return nil
}

func NewMockEmailProvider() *MockEmailProvider {
	return &MockEmailProvider{}
}

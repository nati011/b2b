package sms

type MockProvider struct {
}

func NewMockProvider() MockProvider {
	return MockProvider{}
}
func (m MockProvider) Send(r Request) error {
	return nil
}

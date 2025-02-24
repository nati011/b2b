package email

type Mock struct {
}

func (m *Mock) Send(r Request) error {
	return nil
}

func NewMock() *Mock {
	return &Mock{}
}

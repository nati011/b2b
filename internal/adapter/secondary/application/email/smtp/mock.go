package email

import port "b2b.nati011.github.com/internal/port/application/email"

type Mock struct {
}

func (m *Mock) Send(r port.Request) error {
	return nil
}

func NewMock() *Mock {
	return &Mock{}
}

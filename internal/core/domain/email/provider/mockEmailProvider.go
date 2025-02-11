package email

import emailDTO "b2b.nati011.github.com/internal/core/domain/email/service"

type MockEmailProvider struct {
}

func (m *MockEmailProvider) Send(e emailDTO.Request) (emailDTO.Response, error) {
	return emailDTO.Response{
		Addr:    e.Addr,
		Message: emailDTO.SUCCESS_MESSAGE,
	}, nil
}

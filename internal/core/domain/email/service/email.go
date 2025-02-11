package email

import (
	"errors"

	provider "b2b.nati011.github.com/internal/core/domain/email/provider"
)

type Request struct {
	Addr    string
	Content string
	Header  string
}

type Response struct {
	Addr    string
	Message string
}

// exportable errors
var (
	ErrAddressNotValid = errors.New("oopsy, email is not valid")
	ErrContentEmpty    = errors.New("oopsy, email content is empty")
)

const (
	SUCCESS_MESSAGE = "Ahoy, mail received!"
)

type Emailer interface {
	Send(Request) (Response, error)
}

type EmailService struct {
	provider provider.EmailProvider
}

func NewEmailService(ep provider.EmailProvider) *EmailService {
	return &EmailService{provider: ep}
}

func (e *EmailService) Send(r *Request) (Response, error) {
	response, err := e.provider.Send(
		provider.Request{
			Addr:    r.Addr,
			Header:  r.Header,
			Content: r.Content,
		},
	)
	if err != nil {
		return Response{
			Message: err.Error(),
		}, err
	}
	return Response{
		Addr:    response.Addr,
		Message: SUCCESS_MESSAGE,
	}, nil
}

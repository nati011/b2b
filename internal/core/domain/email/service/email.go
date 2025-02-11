package email

import (
	"errors"
	"regexp"

	provider "b2b.nati011.github.com/internal/core/domain/email/provider"
)

type Request struct {
	From    string
	To      string
	Subject string
	Text    string
}

type Response struct {
	Message string
}

// exportable errors
var (
	ErrSenderAddressNotValid   = errors.New("oopsy, sender email is not valid")
	ErrReceiverAddressNotValid = errors.New("oopsy, receiver email is not valid")
	ErrContentEmpty            = errors.New("oopsy, email content is empty")
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
	isSenderAddrValid := validateEmail(r.From)
	if !isSenderAddrValid {
		return Response{
			Message: ErrSenderAddressNotValid.Error(),
		}, ErrSenderAddressNotValid
	}
	isReceiverAddrValid := validateEmail(r.To)
	if !isReceiverAddrValid {
		return Response{
			Message: ErrReceiverAddressNotValid.Error(),
		}, ErrReceiverAddressNotValid
	}

	isTextNotEmpty := validateMailContent(r.Text)
	if !isTextNotEmpty {
		return Response{
			Message: ErrContentEmpty.Error(),
		}, ErrContentEmpty
	}

	err := e.provider.Send(
		provider.Request{
			From:    r.From,
			To:      r.To,
			Subject: r.Subject,
			Text:    r.Text,
		},
	)
	if err != nil {
		return Response{
			Message: err.Error(),
		}, err
	}
	return Response{
		Message: SUCCESS_MESSAGE,
	}, nil
}

func validateEmail(email string) bool {
	var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func validateMailContent(text string) bool {
	return text != ""
}

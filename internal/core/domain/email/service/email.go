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
	err := validateEmail(r)
	if err != nil {
		return Response{
			Message: err.Error(),
		}, err
	}

	err = e.provider.Send(
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

func validateEmail(r *Request) error {
	isSenderAddrValid := validateEmailAddr(r.From)
	if !isSenderAddrValid {
		return ErrSenderAddressNotValid
	}

	isReceiverAddrValid := validateEmailAddr(r.To)
	if !isReceiverAddrValid {
		return ErrReceiverAddressNotValid
	}

	isTextNotEmpty := validateMailContent(r.Text)
	if !isTextNotEmpty {
		return ErrContentEmpty
	}
	return nil
}

func validateEmailAddr(email string) bool {
	var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func validateMailContent(text string) bool {
	return text != ""
}

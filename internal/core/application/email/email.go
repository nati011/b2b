package email

import (
	"errors"

	smtp "b2b.nati011.github.com/internal/adapter/secondary/application/email/smtp"
	render "b2b.nati011.github.com/internal/core/application/email/render"
)

type SendRequest struct {
	To         string
	Subject    string
	Args       map[string]string
	ExternalId string
}

type GetResponse struct {
	Id         string
	ExternalId string
	Email      Email
}

type Email struct {
	From    string
	To      string
	Subject string
	Text    string
}

var (
	ErrReceiverAddressNotValid = errors.New("oopsy, receiver email is not valid")
	ErrContentEmpty            = errors.New("oopsy, email content is empty")
	ErrUnknown                 = errors.New("oopsy, unknown error")
)

const (
	SUCCESS_MESSAGE = "Ahoy, mail received!"
)

type Emailer interface {
	Send(*SendRequest) error
}

type EmailService struct {
	smtp     smtp.Provider
	renderer render.Renderer
}

func NewEmailService(ep smtp.Provider, r render.Renderer) *EmailService {
	return &EmailService{
		smtp:     ep,
		renderer: r,
	}
}

func (e EmailService) Send(r *SendRequest) error {
	err := validateEmailAddr(r.To)
	if err != nil {
		return err
	}

	renderResponse, err := e.renderer.Create(&render.Request{})
	if err != nil {
		return err
	}
	err = validateMailContent(renderResponse.Text)
	if err != nil {
		return err
	}

	err = e.smtp.Send(
		smtp.Request{
			To:      r.To,
			Subject: r.Subject,
			Text:    renderResponse.Text,
		},
	)
	if err != nil {
		switch err {
		case smtp.ErrSysUnknown:
			return err
		default:
			return ErrUnknown
		}
	}
	return nil
}

func (e EmailService) Get(s string) (GetResponse, error) {
	return GetResponse{}, nil
}

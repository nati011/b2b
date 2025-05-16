package email

import (
	"errors"
	"log"

	render "b2b.nati011.github.com/internal/core/application/render"
	port "b2b.nati011.github.com/internal/port/application/email"
	port_commons "b2b.nati011.github.com/internal/port/commons/db"
)

type SendRequest struct {
	To         string
	Subject    string
	TemplateId int
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

type Provider interface {
	Send(*SendRequest) error
}

type EmailService struct {
	smtp     port.Provider
	renderer render.Provider
}

func NewEmailService(ep port.Provider, r render.Provider) *EmailService {
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

	renderResponse, err := e.renderer.Create(&render.Request{
		TemplateId: r.TemplateId,
		Args:       r.Args,
	})
	if err != nil {
		log.Printf("failed to render email")
		return err
	}

	err = validateMailContent(renderResponse.Text)
	if err != nil {
		return err
	}

	err = e.smtp.Send(
		port.Request{
			To:      r.To,
			Subject: r.Subject,
			Text:    renderResponse.Text,
		},
	)
	if err != nil {
		switch err {
		case port_commons.ErrSysUnknown:
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

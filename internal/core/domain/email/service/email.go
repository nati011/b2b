package email

import (
	"errors"

	smtp "b2b.nati011.github.com/internal/core/domain/email/provider/smtp"
	render "b2b.nati011.github.com/internal/core/domain/email/service/render"
)

type SendRequest struct {
	To         string
	Subject    string
	Args       map[string]string
	ExternalId string
}

type SendResponse struct {
	Message string
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
)

const (
	SUCCESS_MESSAGE = "Ahoy, mail received!"
)

type Emailer interface {
	Send(*SendRequest) (SendResponse, error)
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

func (e EmailService) Send(r *SendRequest) (SendResponse, error) {
	isEmailValid := validateEmailAddr(r.To)
	if !isEmailValid {
		return SendResponse{
			Message: ErrReceiverAddressNotValid.Error(),
		}, ErrReceiverAddressNotValid
	}

	renderResponse, err := e.renderer.Create(&render.Request{})
	if err != nil {
		return SendResponse{
			err.Error(),
		}, err
	}
	isTextValid := validateMailContent(renderResponse.Text)
	if !isTextValid {
		return SendResponse{}, ErrContentEmpty
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
			return SendResponse{}, err
		}
	}

	if err != nil {
		return SendResponse{
			Message: err.Error(),
		}, err
	}

	return SendResponse{
		Message: SUCCESS_MESSAGE,
	}, nil
}

func (e EmailService) Get(s string) (GetResponse, error) {
	return GetResponse{}, nil
}

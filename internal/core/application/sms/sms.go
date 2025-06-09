package sms

import (
	"errors"

	smsProvider "b2b.nati011.github.com/internal/adapter/secondary/application/sms"
)

var (
	ErrPhoneInvalid = errors.New("¯\\_(ツ)_/¯, phone is invalid")
	ErrTextInvalid  = errors.New("¯\\_(ツ)_/¯, text is invalid")
	ErrUnknown      = errors.New("¯\\_(ツ)_/¯, unknown error")
)

const (
	SUCCESS_MESSAGE = "Ahoy!"
)

type Request struct {
	Phone   string
	Content string
}

type Response struct {
	Phone   string
	Message string
}

type Provider interface {
	Send(r *Request) (Response, error)
}

type SMSService struct {
	provider smsProvider.Provider
}

func NewSMSService(sp smsProvider.Provider) *SMSService {
	return &SMSService{
		provider: sp,
	}
}

func (t *SMSService) Send(r *Request) (Response, error) {
	err := validateSMS(r)
	if err != nil {
		return Response{
			Phone:   r.Phone,
			Message: err.Error(),
		}, err
	}

	err = t.provider.Send(smsProvider.Request{
		Phone:   r.Phone,
		Content: r.Content,
	})
	if err != nil {
		return Response{
			Phone:   r.Phone,
			Message: err.Error(),
		}, err
	}

	return Response{
		Phone:   r.Phone,
		Message: SUCCESS_MESSAGE,
	}, nil
}

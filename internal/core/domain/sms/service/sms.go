package sms

import (
	"errors"
	"regexp"
	"strings"

	smsProvider "b2b.nati011.github.com/internal/core/domain/sms/provider"
)

var (
	ErrPhoneInvalid = errors.New("oopsy, phone is invalid")
	ErrTextInvalid  = errors.New("oopsy, text is invalid")
	ErrUnknown      = errors.New("oopsy, unknown error")
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

type Texter interface {
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

func validateSMS(r *Request) error {
	if !validatePhone(r.Phone) {
		return ErrPhoneInvalid
	}
	if !validateContent(r.Content) {
		return ErrTextInvalid
	}
	return nil
}

func validatePhone(p string) bool {
	phoneRegex := regexp.MustCompile(`^\d{10}$`)
	return phoneRegex.MatchString(p)
}

func validateContent(p string) bool {
	return p != "" && (strings.Split(p, "") != nil)
}

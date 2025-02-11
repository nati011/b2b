package email

import "errors"

type Request struct {
	Addr    string
	Content string
	Header  string
}

type Response struct {
	Addr    string
	Message string
}

var (
	ErrAddressNotValid = errors.New("oopsy, email is not valid")
	ErrContentEmpty    = errors.New("oopsy, email content is empty")
)

type Emailer interface {
	SendEmail(Request) (Response, error)
}

type EmailService struct {
}

func NewEmailService() *EmailService {
	return &EmailService{}
}

func (e *EmailService) SendEMail(Request) (Response, error) {
	return Response{}, nil
}

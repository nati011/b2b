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

const (
	SUCCESS_MESSAGE = "Ahoy, mail received!"
)

type Emailer interface {
	Send(Request) (Response, error)
}

type EmailService struct {
}

func NewEmailService() *EmailService {
	return &EmailService{}
}

func (e *EmailService) Send(Request) (Response, error) {
	return Response{}, nil
}

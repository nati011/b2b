package email

import "errors"

var (
	ErrSysUnknown = errors.New("unknown error")
)

type Request struct {
	From    string
	To      string
	Subject string
	Text    string
}

type EmailProvider interface {
	Send(Request) error
}

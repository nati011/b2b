package email

import "errors"

var (
	ErrSysUnknown = errors.New("unknown error")
	ErrFailed     = errors.New("failed to send")
)

type Request struct {
	To      string
	Subject string
	Text    string
}

type Provider interface {
	Send(Request) error
}

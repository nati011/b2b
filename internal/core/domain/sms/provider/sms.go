package sms

type Request struct {
	Phone   string
	Content string
}

type Provider interface {
	Send(Request) error
}

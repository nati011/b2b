package email

type Request struct {
	From    string
	To      string
	Subject string
	Text    string
}

type EmailProvider interface {
	Send(Request) error
}

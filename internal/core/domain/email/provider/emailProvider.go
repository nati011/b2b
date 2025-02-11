package email

type Request struct {
	Addr    string
	Content string
	Header  string
}

type Response struct {
	Addr string
}

type EmailProvider interface {
	Send(Request) (Response, error)
}

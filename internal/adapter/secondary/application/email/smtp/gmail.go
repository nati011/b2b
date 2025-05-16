package email

import (
	"log"
	"net/smtp"

	port "b2b.nati011.github.com/internal/port/application/email"
)

type Gmail struct {
	sender   string
	smtpPort string
}

func NewGmail(sender string, smtpPort string) port.Provider {
	return &Gmail{
		sender:   sender,
		smtpPort: smtpPort,
	}
}

func (m Gmail) Send(r port.Request) error {

	c, err := smtp.Dial(m.smtpPort)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	if err := c.Mail(m.sender); err != nil {
		log.Fatal(err)
	}
	if err := c.Rcpt(r.To); err != nil {
		log.Fatal(err)
	}

	wc, err := c.Data()
	if err != nil {
		log.Fatal(err)
	}

	msg := []byte("To:" + r.To + "\r\n" +
		"From:" + m.sender + "\r\n" +
		"Subject:" + r.Subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n" +
		"\r\n" +
		"" + r.Text + "\r\n")

	_, err = wc.Write(msg)
	if err != nil {
		log.Fatal(err)
	}

	err = wc.Close()
	if err != nil {
		log.Fatal(err)
	}

	if err := c.Quit(); err != nil {
		log.Fatal(err)
	}

	return nil
}

package email

import (
	"crypto/tls"
	"log"
	"net/smtp"

	port "b2b.nati011.github.com/internal/port/application/email"
)

type Gmail struct {
	sender   string
	smtpPort string
	password string
}

func NewGmail(sender, smtpPort, password string) port.Provider {
	return &Gmail{
		sender:   sender,
		smtpPort: smtpPort,
		password: password,
	}
}

func (m Gmail) Send(r port.Request) error {
	// Set up authentication information.
	auth := smtp.PlainAuth("", m.sender, m.password, "smtp.gmail.com")

	// Connect to the SMTP server.
	c, err := smtp.Dial(m.smtpPort)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	// Start TLS
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         "smtp.gmail.com",
	}
	if err = c.StartTLS(tlsconfig); err != nil {
		log.Fatal(err)
	}

	// Authenticate
	if err = c.Auth(auth); err != nil {
		log.Fatal(err)
	}

	// Set the sender and recipient
	if err = c.Mail(m.sender); err != nil {
		log.Fatal(err)
	}
	if err = c.Rcpt(r.To); err != nil {
		log.Fatal(err)
	}

	// Send the email body
	wc, err := c.Data()
	if err != nil {
		log.Fatal(err)
	}

	msg := []byte("To: " + r.To + "\r\n" +
		"From: " + m.sender + "\r\n" +
		"Subject: " + r.Subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n" +
		"\r\n" +
		r.Text + "\r\n")

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

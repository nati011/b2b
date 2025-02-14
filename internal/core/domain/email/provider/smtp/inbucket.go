package email

import (
	"log"
	"net/smtp"
)

type Inbucket struct {
	sender   string
	smtpPort string
}

func NewInbucket(sender string, smtpPort string) Provider {
	return &Inbucket{
		sender:   sender,
		smtpPort: smtpPort,
	}
}

func (m Inbucket) Send(r Request) error {
	// Connect to the remote SMTP server.
	c, err := smtp.Dial(m.smtpPort)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close() // Ensure the connection is closed

	// Set the sender and recipient first
	if err := c.Mail(m.sender); err != nil {
		log.Fatal(err)
	}
	if err := c.Rcpt(r.To); err != nil {
		log.Fatal(err)
	}

	// Send the email body.
	wc, err := c.Data()
	if err != nil {
		log.Fatal(err)
	}

	// Prepare the email headers and body
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

	// Send the QUIT command and close the connection.
	if err := c.Quit(); err != nil {
		log.Fatal(err)
	}

	return nil
}

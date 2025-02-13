package email

import (
	"os"
	"testing"
)

const (
	VALID_EMAIL_RECEPIENT = "ruthtirusew944@mailpit.com"
	INVALID_EMAIL_ADDR    = "ruthtirusew944"
	VALID_SUBJECT         = "test"
	VALID_TEXT            = "test"
	INVALID_TEXT          = ""
)

var service *EmailService

func Test_SendEmail_happyPath(t *testing.T) {
}

func Test_SendEmail_unhappyPath(t *testing.T) {

	t.Run("invalidRecepientEmail", func(t *testing.T) {

	})

	t.Run("emptyContent", func(t *testing.T) {
	})
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	// service = NewEmailService(smtp.NewMock(), db.NewMock(), er.mock)
}

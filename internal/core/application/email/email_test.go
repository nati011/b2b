package email

import (
	"os"
	"testing"

	smtp "b2b.nati011.github.com/internal/adapter/secondary/application/email/smtp"
	render "b2b.nati011.github.com/internal/core/application/render"
)

const (
	VALID_EMAIL_RECEPIENT = "ruthtirusew944@mailpit.com"
	INVALID_EMAIL_ADDR    = "ruthtirusew944"
	VALID_SUBJECT         = "test"
	VALID_TEXT            = "test"
	INVALID_TEXT          = ""
)

var service Emailer

func Test_SendEmail_happyPath(t *testing.T) {
	in := SendRequest{
		To:      VALID_EMAIL_RECEPIENT,
		Subject: VALID_SUBJECT,
		Args: map[string]string{
			"test": "test",
		},
	}

	err := service.Send(&in)
	if err != nil {
		t.Fatalf("Failed to send err: %v", err)
	}
}

func Test_SendEmail_unhappyPath(t *testing.T) {
	t.Run("invalidRecepientEmail", func(t *testing.T) {
		in := SendRequest{
			To:      INVALID_EMAIL_ADDR,
			Subject: VALID_SUBJECT,
			Args: map[string]string{
				"test": "test",
			},
		}

		err := service.Send(&in)
		wantErr := ErrReceiverAddressNotValid
		if err != wantErr {
			t.Errorf("Expected err: %v Got: %v", wantErr, err)
		}
	})
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	service = NewEmailService(smtp.NewMock(), render.NewMock())
}

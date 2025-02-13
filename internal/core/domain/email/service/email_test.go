package email

import (
	"os"
	"testing"

	smtp "b2b.nati011.github.com/internal/core/domain/email/provider/smtp"
	render "b2b.nati011.github.com/internal/core/domain/email/service/render"

	db "b2b.nati011.github.com/internal/core/domain/email/provider/db/email"
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
	want := SendResponse{
		Message: SUCCESS_MESSAGE,
	}

	got, err := service.Send(&in)
	if err != nil {
		t.Errorf("Failed to send email err: %q", err)
	}
	if got != want {
		t.Errorf("Expecetd: %q Got: %q", want, got)
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
		want := SendResponse{
			Message: ErrReceiverAddressNotValid.Error(),
		}
		got, err := service.Send(&in)
		if err != ErrReceiverAddressNotValid {
			t.Errorf("Failed to send email %q", err)
		}
		if got != want {
			t.Errorf("Expected: %q Got: %q", want, got)
		}
	})
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	service = NewEmailService(smtp.NewMock(), render.NewMock(), db.NewMock())
}

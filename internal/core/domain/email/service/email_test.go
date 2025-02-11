package email

import (
	"os"
	"testing"

	provider "b2b.nati011.github.com/internal/core/domain/email/provider"
)

const (
	VALID_EMAIL_RECEPIENT = "ruthtirusew944@mailpit.com"
	VALID_EMAIL_SENDER    = "ruthtirusew388@gmail.com"
	INVALID_EMAIL_ADDR    = "ruthtirusew944"
	VALID_SUBJECT         = "test"
	VALID_TEXT            = "test"
	INVALID_TEXT          = ""
)

var service *EmailService

func Test_SendEmail_happyPath(t *testing.T) {
	in := Request{
		From:    VALID_EMAIL_SENDER,
		To:      VALID_EMAIL_RECEPIENT,
		Subject: VALID_SUBJECT,
		Text:    VALID_TEXT,
	}
	want := Response{
		Message: SUCCESS_MESSAGE,
	}

	got, _ := service.Send(&in)
	// if err != nil {
	// 	t.Errorf("Failed to send email err: %v", err)
	// }
	if got != want {
		t.Errorf("Expected: %v, Got: %v", want, got)
	}
}

func Test_SendEmail_unhappyPath(t *testing.T) {
	t.Run("invalidSenderEmail", func(t *testing.T) {
		in := Request{
			From:    INVALID_EMAIL_ADDR,
			To:      VALID_EMAIL_RECEPIENT,
			Subject: VALID_SUBJECT,
			Text:    VALID_TEXT,
		}
		want := Response{
			Message: ErrSenderAddressNotValid.Error(),
		}

		got, _ := service.Send(&in)
		// if err != nil {
		// 	t.Errorf("Failed to send email err: %v", err)
		// }
		if got != want {
			t.Errorf("Expected: %v, Got: %v", want, got)
		}
	})

	t.Run("invalidRecepientEmail", func(t *testing.T) {
		in := Request{
			From:    VALID_EMAIL_SENDER,
			To:      INVALID_EMAIL_ADDR,
			Subject: VALID_SUBJECT,
			Text:    VALID_TEXT,
		}
		want := Response{
			Message: ErrReceiverAddressNotValid.Error(),
		}

		got, _ := service.Send(&in)
		// if err != nil {
		// 	t.Errorf("Failed to send email err: %v", err)
		// }
		if got != want {
			t.Errorf("Expected: %v, Got: %v", want, got)
		}
	})

	t.Run("emptyContent", func(t *testing.T) {
		in := Request{
			From:    VALID_EMAIL_SENDER,
			To:      VALID_EMAIL_RECEPIENT,
			Subject: VALID_SUBJECT,
			Text:    INVALID_TEXT,
		}
		want := Response{
			Message: ErrContentEmpty.Error(),
		}
		got, _ := service.Send(&in)
		// if err != nil {
		// 	t.Errorf("Failed to send email err: %v", err)
		// }
		if got != want {
			t.Errorf("Expected: %v, Got: %v", want, got)
		}
	})
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	service = NewEmailService(provider.NewMockEmailProvider())
}

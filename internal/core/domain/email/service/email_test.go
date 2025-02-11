package email

import (
	"os"
	"testing"
)

const (
	VALID_EMAIL_ADDR   = "ruthtirusew944@mailpit.com"
	INVALID_EMAIL_ADDR = "ruthtirusew944"
	VALID_HEADER       = "test"
	VALID_CONTENT      = "test"
)

var service *EmailService

func Test_SendEmail_happyPath(t *testing.T) {
	in := Request{
		Addr:    VALID_EMAIL_ADDR,
		Header:  VALID_HEADER,
		Content: VALID_CONTENT,
	}

	want := Response{
		Addr:    VALID_EMAIL_ADDR,
		Message: SUCCESS_MESSAGE,
	}

	got, err := service.Send(in)
	if err != nil {
		t.Errorf("Failed to send email err: %v", err)
	}
	if got != want {
		t.Errorf("Expected: %v, Got: %v", want, got)
	}
}

func Test_SendEmail_unhappyPath(t *testing.T) {
	t.Run("invalidEmail", func(t *testing.T) {
		in := Request{
			Addr:    INVALID_EMAIL_ADDR,
			Header:  VALID_HEADER,
			Content: VALID_CONTENT,
		}
		want := Response{
			Addr:    INVALID_EMAIL_ADDR,
			Message: ErrAddressNotValid.Error(),
		}
		got, err := service.Send(in)
		if err != nil {
			t.Errorf("Failed to send email err: %v", err)
		}
		if got != want {
			t.Errorf("Expected: %v, Got: %v", want, got)
		}

	})

	t.Run("emptyContent", func(t *testing.T) {
		in := Request{
			Addr:    VALID_EMAIL_ADDR,
			Header:  VALID_HEADER,
			Content: VALID_CONTENT,
		}
		want := Response{
			Addr:    VALID_EMAIL_ADDR,
			Message: ErrContentEmpty.Error(),
		}
		got, err := service.Send(in)
		if err != nil {
			t.Errorf("Failed to send email err: %v", err)
		}
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
	service = NewEmailService()
}

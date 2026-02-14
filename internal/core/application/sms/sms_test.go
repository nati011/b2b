package sms

import (
	"os"
	"testing"

	smsProvider "b2b.nati011.github.com/internal/adapter/secondary/application/sms"
)

const (
	VALID_PHONE     = "0949184879"
	INVALID_PHONE   = "011"
	VALID_CONTENT   = "test"
	INVALID_CONTENT = " "
	EMPTY_CONTENT   = ""
)

var service *SMSService

func Test_Sendsms_happyPath(t *testing.T) {
	in := Request{
		Phone:   VALID_PHONE,
		Content: VALID_CONTENT,
	}
	want := Response{
		Phone:   VALID_PHONE,
		Message: SUCCESS_MESSAGE,
	}
	got, err := service.Send(&in)
	if err != nil {
		t.Errorf("Failed to send sms err: %v", err)
	}
	if want != got {
		t.Errorf("Expected: %v, Got: %v", want, got)
	}
}

func Test_Sendsms_unhappyPath(t *testing.T) {
	t.Run("invalidPhone", func(t *testing.T) {
		in := Request{
			Phone:   INVALID_PHONE,
			Content: VALID_CONTENT,
		}
		want := Response{
			Phone:   INVALID_PHONE,
			Message: ErrPhoneInvalid.Error(),
		}
		got, err := service.Send(&in)
		if err != ErrPhoneInvalid {
			t.Errorf("Failed to send sms err: %v", err)
		}
		if want != got {
			t.Errorf("Expected: %v, Got: %v", want, got)
		}
	})

	t.Run("emptyText", func(t *testing.T) {
		in := Request{
			Phone:   VALID_PHONE,
			Content: EMPTY_CONTENT,
		}
		want := Response{
			Phone:   VALID_PHONE,
			Message: ErrTextInvalid.Error(),
		}
		got, err := service.Send(&in)
		if err != ErrTextInvalid {
			t.Errorf("Failed to send sms err: %v", err)
		}
		if want != got {
			t.Errorf("Expected: %v, Got: %v", want, got)
		}
	})

	t.Run("invalidText", func(t *testing.T) {
		in := Request{
			Phone:   VALID_PHONE,
			Content: INVALID_CONTENT,
		}
		want := Response{
			Phone:   VALID_PHONE,
			Message: SUCCESS_MESSAGE,
		}
		got, err := service.Send(&in)
		if err != nil {
			t.Errorf("Failed to send sms err: %v", err)
		}
		if want != got {
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
	service = NewSMSService(smsProvider.NewMockProvider())
}

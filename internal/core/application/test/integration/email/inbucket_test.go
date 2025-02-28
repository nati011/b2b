package integration

import (
	"context"
	"os"
	"testing"

	emailProvider "b2b.nati011.github.com/pkg/email/provider/smtp"
	email "b2b.nati011.github.com/pkg/email/service"
	"github.com/testcontainers/testcontainers-go/modules/inbucket"
)

var inbucketContainer *inbucket.InbucketContainer
var container email.TestContainer

const (
	VALID_EMAIL_RECEPIENT = "ruthtirusew944@mailpit.com"
	INVALID_EMAIL_ADDR    = "ruthtirusew944"
	VALID_SUBJECT         = "test"
	VALID_TEXT            = "test"
	INVALID_TEXT          = ""
)

func Test_Send_Email(t *testing.T) {
	in := email.SendRequest{
		To:      VALID_EMAIL_RECEPIENT,
		Subject: VALID_SUBJECT,
		Args: map[string]string{
			"test": "test",
		},
	}
	want := email.SendResponse{
		Message: email.SUCCESS_MESSAGE,
	}

	got, err := container.Emailer.Send(&in)
	if err != nil {
		t.Errorf("Failed to send email err: %q", err)
	}
	if got != want {
		t.Errorf("Expected: %q Got: %q", want, got)
	}
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutDown()
	os.Exit(code)
}

func setup() {
	var err error
	ctx := context.Background()
	inbucketContainer, err = RunContainer(ctx)
	if err != nil {
		panic(err)
	}

	smtpPort, err := inbucketContainer.SmtpConnection(ctx)
	if err != nil {
		panic(err)
	}

	container = *email.NewTestContainer(emailProvider.NewInbucket(
		"test@gmail.com",
		smtpPort,
	))

}

const (
	INBUCKET_VERSION = "inbucket/inbucket:sha-2d409bb"
)

func RunContainer(ctx context.Context) (*inbucket.InbucketContainer, error) {
	return inbucket.Run(ctx, INBUCKET_VERSION)
}

func shutDown() {
	ctx := context.Background()
	err := inbucketContainer.Terminate(ctx)
	if err != nil {
		panic(err)
	}
}

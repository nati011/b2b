package integration

import (
	"context"
	"os"
	"testing"

	emailProvider "b2b.nati011.github.com/internal/adapter/secondary/application/email/smtp"
	email "b2b.nati011.github.com/internal/core/application/email"
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
	err := container.Emailer.Send(&in)
	if err != nil {
		t.Fatalf("Failed to send email err: %q", err)
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

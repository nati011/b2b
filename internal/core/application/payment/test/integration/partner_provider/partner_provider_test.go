package payment_provider

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
}

func Test_Init_Payment_On_Checkout(t *testing.T) {
}

func Test_Handle_Webhook_Payment_Notification(t *testing.T) {
}

func Test_Payment_Verification(t *testing.T) {
}

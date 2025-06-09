package order_confirmation

import "testing"

func TestMain(t *testing.M) {
	setup()
}

func setup() {
}

func teardown() {
}

func Test_confirm_order_happyPath(t *testing.T) {
	t.Run("changeOrderConfirmationStatusToConfirmed", func(t *testing.T) {
		t.Cleanup(teardown)
	})
}

func Test_confirm_order_unhappyPath(t *testing.T) {
	t.Run("orderNotFound", func(t *testing.T) {
		t.Cleanup(teardown)
	})

	t.Run("orderAlreadyConfirmed", func(t *testing.T) {
		t.Cleanup(teardown)
	})
}

func Test_reject_order_happyPath(t *testing.T) {
	t.Run("changeOrderConfirmationStatusToRejected", func(t *testing.T) {
		t.Cleanup(teardown)
	})
}

func Test_reject_order_unhappyPath(t *testing.T) {
	t.Run("orderNotFound", func(t *testing.T) {
		t.Cleanup(teardown)
	})

	t.Run("orderAlreadyRejected", func(t *testing.T) {
		t.Cleanup(teardown)
	})
}

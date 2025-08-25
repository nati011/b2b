package integration

import "testing"

func TestMain(m *testing.M) {

}

func setup() {

}

func teardown() {

}

func Test_Checkout_happypath(t *testing.T) {
	t.Run("generate_checkout_url_upon_placement_digital", func(t *testing.T) {
		t.Cleanup(teardown)
	})

	t.Run("checkout_upon_placement_manual", func(t *testing.T) {
		t.Cleanup(teardown)
	})

	t.Run("generate_checkout_url_upon_init_payment", func(t *testing.T) {
		t.Cleanup(teardown)
	})
}

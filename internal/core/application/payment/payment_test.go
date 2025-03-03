package payment

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

func Test_Init(t *testing.T) {
}

func Test_Cancel(t *testing.T) {
}

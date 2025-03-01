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

func Test_Create_Payment_Option(t *testing.T) {
}

func Test_Activate_Payment_Option(t *testing.T) {
}

func Test_Deactivate_Payment_Option(t *testing.T) {
}

func Test_Get_All_Payment_Options(t *testing.T) {
}

func Test_Get_Active_Payment_Options(t *testing.T) {
}

func Test_Checkout(t *testing.T) {
}

func Test_Payment_Confirmation(t *testing.T) {
}

func Test_Get_Payment_Status(t *testing.T) {
}

func Test_Save_Transaction(t *testing.T) {
}

func Test_Get_Transaction_History(t *testing.T) {
}

func Test_Get_Transaction_History_By_User(t *testing.T) {
}

func Test_Get_Transaction_History_By_Status(t *testing.T) {
}

func Test_GetAll_Transaction_History(t *testing.T) {
}

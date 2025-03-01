package payment

const (
	CALLBACK_URL = "/payment_call_back"
)

type CreatePaymentOptionRequest struct {
	Name string
	Icon string
}

type GetPaymentOptionResponse struct {
	Id     int
	Name   string
	Icon   string
	Status string
}

type GetAllPaymentOptionsResponse struct {
	List []GetPaymentOptionResponse
}

type CheckoutRequest struct {
	User_Id     int
	PhoneNumber string
	Amount      int64
	// Return_url   string
	// Callback_url string
}

type CheckoutResponse struct {
	Checkout_url string
}

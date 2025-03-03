package payment

type CheckoutRequest struct {
	User_Id     int
	PhoneNumber string
	Amount      int64
	// Return_url   string
	// Callback_url string
}

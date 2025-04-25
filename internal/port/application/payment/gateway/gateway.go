package payment

import "errors"

var (
	ErrUnknown            = errors.New("oopsy, Unknown error has occured")
	ErrInvalidTransaction = errors.New("oopsy, the transaction is invalid")
)

type InitiateRequest struct {
	Amount         float64
	BaseUrl        string
	PartnerUrl     string
	PartnerSecret  string
	ReturnUrl      string
	TransactionRef string
}

type InitatePaymentResponse struct {
	CheckoutUrl string
}

type VerifyRequest struct {
	PartnerUrl     string
	TransactionRef string
	PartnerSecret  string
}

type VerifyResponse struct {
	Status string
}

type Provider interface {
	Initiate(req InitiateRequest) (InitatePaymentResponse, error)
	Verify(VerifyRequest) (VerifyResponse, error)
}

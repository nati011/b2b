package payment

import "errors"

var (
	ErrUnknown            = errors.New("oopsy, Unknown error has occured")
	ErrInvalidTransaction = errors.New("oopsy, the transaction is invalid")
)

type InitiateRequest struct {
	Amount         float64 `json:"amount"`
	CallbackUrl    string  `json:"callback_url"`
	PartnerUrl     string  `json:"partner_url"`
	PartnerSecret  string  `json:"partner_secret"`
	ReturnUrl      string  `json:"return_url"`
	TransactionRef int     `json:"tx_ref"`
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

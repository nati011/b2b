package payment

import "errors"

var (
	ErrUnknown            = errors.New("oopsy, Unknown error has occured")
	ErrInvalidTransaction = errors.New("oopsy, error invalid transaction")
)

type InitiateRequest struct {
	Amount         float64 `json:"amount"`
	TransactionRef string  `json:"tx_ref"`
	ReturnUrl      string  `json:"return_url"`
	PartnerUrl     string  `json:"partner_url"`
	PartnerSecret  string  `json:"partner_secret"`
	BaseUrl        string
}

type VerificationRequest struct {
	TransactionRef string `json:"tx_ref"`
	PartnerUrl     string `json:"partner_url"`
	PartnerSecret  string `json:"partner_secret"`
}

type Provider interface {
	Initiate(req InitiateRequest) (string, error)
	Verify(req VerificationRequest) (bool, error)
}

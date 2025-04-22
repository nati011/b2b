package payment

import "errors"

var (
	ErrUnknown = errors.New("oopsy, Unknown error has occured")
)

type InitiateRequest struct {
	Amount         float64 `json:"amount"`
	TransactionRef int     `json:"tx_ref"`
	ReturnUrl      string  `json:"return_url"`
	PartnerUrl     string  `json:"partner_url"`
	PartnerSecret  string  `json:"partner_secret"`
}

type Provider interface {
	Initiate(req InitiateRequest) (string, error)
	Verify()
}

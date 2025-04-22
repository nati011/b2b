package gateway

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	port "b2b.nati011.github.com/internal/port/application/payment/gateway"
)

var (
	ErrPartnerUrlIsMandatory = errors.New("oopsy, partner url is mandatory")
)

type InitatePaymentChapaResponseData struct {
	CheckoutUrl string `json:"checkout_url"`
}

type InitatePaymentChapaResponse struct {
	Message string                          `json:"message"`
	Status  string                          `json:"status"`
	Data    InitatePaymentChapaResponseData `json:"data"`
}

type Chapa struct {
}

func NewChapa() port.Provider {
	return &Chapa{}
}

func (t Chapa) Initiate(request port.InitiateRequest) (string, error) {
	url := request.PartnerUrl
	var response InitatePaymentChapaResponse

	payload, err := json.Marshal(request)
	if err != nil {
		return "", port.ErrUnknown
	}

	client := &http.Client{}
	body := bytes.NewReader(payload)
	req, err := http.NewRequest("POST", url, body)

	if err != nil {
		return "", port.ErrUnknown
	}

	req.Header.Add("Authorization", request.PartnerSecret)
	res, err := client.Do(req)

	if err != nil {
		return "", port.ErrUnknown
	}

	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return "", port.ErrUnknown
	}

	if err := json.Unmarshal(resBody, &response); err != nil {
		return "", port.ErrUnknown
	}

	switch response.Status {
	case "success":
		return response.Data.CheckoutUrl, nil
	default:
		return "", port.ErrUnknown
	}

}

func (t Chapa) Verify() {}

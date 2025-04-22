package gateway

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

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

type InitatePaymentChapaRequest struct {
	Amount         float64 `json:"amount"`
	CallbackUrl    string  `json:"callback_url"`
	ReturnUrl      string  `json:"return_url"`
	TransactionRef int     `json:"tx_ref"`
}

type VerifyPaymentChapaResponseData struct {
	Amount         float64 `json:"amount"`
	Status         string  `json:"status"`
	TransactionRef string  `json:"tx_ref"`
}

type VerifyPaymentChapaResponse struct {
	Message string                          `json:"message"`
	Status  string                          `json:"status"`
	Data    InitatePaymentChapaResponseData `json:"data"`
}

type Chapa struct {
}

func NewChapa() port.Provider {
	return &Chapa{}
}

func (t Chapa) Initiate(request port.InitiateRequest) (port.InitatePaymentResponse, error) {
	var response InitatePaymentChapaResponse

	callback_url := fmt.Sprintf("%v/chapa/%v", request.PartnerUrl, strconv.Itoa(request.TransactionRef))

	requestBody := InitatePaymentChapaRequest{
		Amount:         request.Amount,
		TransactionRef: request.TransactionRef,
		CallbackUrl:    callback_url,
		ReturnUrl:      request.ReturnUrl,
	}

	payload, err := json.Marshal(requestBody)
	if err != nil {
		return port.InitatePaymentResponse{}, port.ErrUnknown
	}

	client := &http.Client{}
	body := bytes.NewReader(payload)
	req, err := http.NewRequest("POST", request.PartnerUrl, body)

	if err != nil {
		return port.InitatePaymentResponse{}, port.ErrUnknown
	}

	req.Header.Add("Authorization", request.PartnerSecret)
	res, err := client.Do(req)

	if err != nil {
		return port.InitatePaymentResponse{}, port.ErrUnknown
	}

	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return port.InitatePaymentResponse{}, port.ErrUnknown
	}

	if err := json.Unmarshal(resBody, &response); err != nil {
		return port.InitatePaymentResponse{}, port.ErrUnknown
	}

	switch response.Status {
	case "success":
		return port.InitatePaymentResponse{
			CheckoutUrl: response.Data.CheckoutUrl}, nil
	default:
		return port.InitatePaymentResponse{}, port.ErrUnknown
	}

}

func (t Chapa) Verify(request port.VerifyRequest) (port.VerifyResponse, error) {
	var response VerifyPaymentChapaResponse

	verification_url := fmt.Sprintf("%v/transaction/verify/%v", request.PartnerUrl, strconv.Itoa(request.TransactionRef))

	client := &http.Client{}
	req, err := http.NewRequest("GET", verification_url, nil)

	if err != nil {
		return port.VerifyResponse{}, port.ErrUnknown
	}

	req.Header.Add("Authorization", request.PartnerSecret)
	res, err := client.Do(req)

	if err != nil {
		return port.VerifyResponse{}, port.ErrUnknown
	}

	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return port.VerifyResponse{}, port.ErrUnknown
	}

	if err := json.Unmarshal(resBody, &response); err != nil {
		return port.VerifyResponse{}, port.ErrUnknown
	}

	switch response.Status {
	case "success":
		return port.VerifyResponse{
			Status: response.Status,
		}, nil
	case "failed":
		return port.VerifyResponse{}, port.ErrInvalidTransaction
	default:
		return port.VerifyResponse{}, port.ErrUnknown
	}
}

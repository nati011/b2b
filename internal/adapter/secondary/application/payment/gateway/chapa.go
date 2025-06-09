package gateway

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	port "b2b.nati011.github.com/internal/port/application/payment/gateway"
)

var (
	ErrPartnerUrlIsMandatory = errors.New("¯\\_(ツ)_/¯, partner url is mandatory")
)

type InitatePaymentChapaResponseData struct {
	CheckoutUrl string `json:"checkout_url"`
}

type InitatePaymentChapaResponse struct {
	Message string                          `json:"message"`
	Status  string                          `json:"status"`
	Data    InitatePaymentChapaResponseData `json:"data"`
}

type VerifyResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
	Data    struct {
		FirstName     string  `json:"first_name"`
		LastName      string  `json:"last_name"`
		Email         string  `json:"email"`
		Currency      string  `json:"currency"`
		Amount        float64 `json:"amount"`
		Charge        float64 `json:"charge"`
		Mode          string  `json:"mode"`
		Method        string  `json:"method"`
		Type          string  `json:"type"`
		Status        string  `json:"status"`
		Reference     string  `json:"reference"`
		TxRef         string  `json:"tx_ref"`
		Customization struct {
			Title       string  `json:"title"`
			Description string  `json:"description"`
			Logo        *string `json:"logo"`
		} `json:"customization"`
		Meta      interface{} `json:"meta"`
		CreatedAt time.Time   `json:"created_at"`
		UpdatedAt time.Time   `json:"updated_at"`
	} `json:"data"`
}

type Chapa struct {
}

func NewChapa() port.Provider {
	return &Chapa{}
}

func (t Chapa) Initiate(request port.InitiateRequest) (string, error) {
	initalization_url := fmt.Sprintf("%v/v1/transaction/initialize", request.PartnerUrl)
	var response InitatePaymentChapaResponse
	var logs any
	payload, err := json.Marshal(request)
	if err != nil {
		return "", port.ErrUnknown
	}
	log.Printf("%v/v1/transaction/initialize", request.PartnerUrl)
	client := &http.Client{}
	body := bytes.NewReader(payload)
	req, err := http.NewRequest("POST", initalization_url, body)

	if err != nil {
		return "", port.ErrUnknown
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", request.PartnerSecret))
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)

	if err != nil {
		log.Printf("Error while checkingout: %v", err.Error())
		return "", port.ErrUnknown
	}

	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("Error while checkingout: %v", err.Error())
		return "", port.ErrUnknown
	}
	if err := json.Unmarshal(resBody, &logs); err != nil {
		log.Printf("Error while checkingout: %v", err.Error())
		return "", port.ErrUnknown
	}
	log.Printf("Gateway response: %v", logs)
	if err := json.Unmarshal(resBody, &response); err != nil {
		log.Printf("Error while checkingout: %v", err.Error())
		return "", port.ErrUnknown
	}

	switch response.Status {
	case "success":
		return response.Data.CheckoutUrl, nil
	default:
		return "", port.ErrUnknown
	}
}

func (t Chapa) Verify(request port.VerificationRequest) (bool, error) {
	var response VerifyResponse

	verification_url := fmt.Sprintf("%v/v1/transaction/verify/%v", request.PartnerUrl, request.TransactionRef)

	client := &http.Client{}
	req, err := http.NewRequest("GET", verification_url, nil)

	if err != nil {
		return false, port.ErrUnknown
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", request.PartnerSecret))
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)

	if err != nil {
		return false, port.ErrUnknown
	}

	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return false, port.ErrUnknown
	}

	if err := json.Unmarshal(resBody, &response); err != nil {
		return false, port.ErrUnknown
	}

	switch response.Data.Status {
	case "success":
		return true, nil
	default:
		return false, nil
	}
}

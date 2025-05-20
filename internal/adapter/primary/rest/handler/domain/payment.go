package domain

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"b2b.nati011.github.com/internal/adapter/primary/rest/handler"
	util "b2b.nati011.github.com/internal/adapter/primary/rest/handler/util"
	application_core "b2b.nati011.github.com/internal/core/application"
	domain_core "b2b.nati011.github.com/internal/core/domain"
	"b2b.nati011.github.com/internal/core/domain/payment_verification"
)

var (
	ErrUnknownPaymentCommand = errors.New("unknown user command")
)

type CheckoutRequest struct {
	OrderId          int `json:"order_id"`
	PaymentPartnerId int `json:"payment_partner_id"`
}

type Payment struct {
	service payment_verification.Provider
}

func InitPayment() {
	handler.Register(new(Payment))
}

func (p *Payment) Init(authMiddleWare *util.AuthMiddleware, applicationServices *application_core.Container, domainService *domain_core.Container) error {
	p.service = domainService.PaymentVerificationService
	return nil
}

func (p *Payment) Routes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/payment/webhook/{gateway_id}/{tx_ref}", p.CallbackGETHandler)
	mux.HandleFunc("POST /api/v1/payment/webhook/{gateway_id}", p.CallbackPOSTHandler)
}

func (p *Payment) CallbackGETHandler(w http.ResponseWriter, r *http.Request) {
	typedParamGatewayId, err := util.GetPathParam(r, 5)
	if err != nil {
		util.RequestErrorResponse(w, err)
		return
	}
	typedParamTxRef, err := util.GetStringPathParam(r, 6)
	if err != nil {
		util.RequestErrorResponse(w, err)
		return
	}
	p.service.Callback(r.Context(), typedParamGatewayId, typedParamTxRef)
}

type CallbackPOSTHandlerBody struct {
	TxRef string `json:"tx_ref"`
}

func (p *Payment) CallbackPOSTHandler(w http.ResponseWriter, r *http.Request) {
	typedParamGatewayId, err := util.GetPathParam(r, 5)
	if err != nil {
		util.RequestErrorResponse(w, err)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		util.RequestErrorResponse(w, err)
		return
	}
	defer r.Body.Close()

	var requestBody CallbackPOSTHandlerBody
	if err := json.Unmarshal(body, &requestBody); err != nil {
		util.RequestErrorResponse(w, err)
		return
	}
	p.service.Callback(r.Context(), typedParamGatewayId, requestBody.TxRef)
}

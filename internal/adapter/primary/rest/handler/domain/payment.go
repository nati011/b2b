package domain

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"b2b.nati011.github.com/internal/adapter/primary/rest/handler"
	util "b2b.nati011.github.com/internal/adapter/primary/rest/handler/util"
	application_core "b2b.nati011.github.com/internal/core/application"
	"b2b.nati011.github.com/internal/core/application/middleware"
	"b2b.nati011.github.com/internal/core/application/resource"
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
	authMiddleware middleware.Auth
	service        payment_verification.Provider
}

func InitPayment() {
	handler.Register(new(Payment))

	handler.RegisterResource(resource.CreateRequest{
		Name:     "payment_webook_callback_var_1",
		Action:   "ALL",
		Resource: "/api/v1/payment/webhook/{gateway_id}/{tx_ref}",
	})
	handler.RegisterResource(resource.CreateRequest{
		Name:     "payment_webook_callback_var_2",
		Action:   "ALL",
		Resource: "/api/v1/payment/webhook/{param}",
	})
	handler.RegisterResource(resource.CreateRequest{
		Name:     "payment_verify",
		Action:   "ALL",
		Resource: "/api/v1/payment/verify",
	})
	handler.RegisterResource(resource.CreateRequest{
		Name:     "payment_confirm",
		Action:   "ALL",
		Resource: "/api/v1/payment/confirm/{param}",
	})
}

func (p *Payment) Init(authMiddleWare *middleware.Auth, applicationServices *application_core.Container, domainService *domain_core.Container) error {
	p.service = domainService.PaymentVerificationService
	p.authMiddleware = *applicationServices.AuthMiddleware
	return nil
}

func (p *Payment) Routes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/payment/webhook/{tx_ref}", func(w http.ResponseWriter, r *http.Request) {
		p.authMiddleware.RequireNoAuthentication(http.HandlerFunc(p.CallbackGETHandler)).ServeHTTP(w, r)
	})

	mux.HandleFunc("POST /api/v1/payment/webhook/{partnerId}", func(w http.ResponseWriter, r *http.Request) {
		p.authMiddleware.RequireNoAuthentication(http.HandlerFunc(p.CallbackPOSTHandler)).ServeHTTP(w, r)
	})

	mux.HandleFunc("POST /api/v1/payment/webhook", func(w http.ResponseWriter, r *http.Request) {
		p.authMiddleware.RequireNoAuthentication(http.HandlerFunc(p.CallbackPOSTHandler)).ServeHTTP(w, r)
	})

	mux.HandleFunc("POST /api/v1/payment/verify", func(w http.ResponseWriter, r *http.Request) {
		p.authMiddleware.RequireAuthentication(http.HandlerFunc(p.VerifyPaymentHandler)).ServeHTTP(w, r)
	})

	mux.HandleFunc("POST /api/v1/payment/confirm/{tx_ref}", func(w http.ResponseWriter, r *http.Request) {
		p.authMiddleware.RequireAuthentication(http.HandlerFunc(p.ConfirmPaymentHandler)).ServeHTTP(w, r)
	})
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

func (p *Payment) VerifyPaymentHandler(w http.ResponseWriter, r *http.Request) {
	const ParamTxRef = "tx_ref"
	paramValues := r.URL.Query()
	paramTxRefValue := paramValues.Get(ParamTxRef)
	if paramTxRefValue != "" {
		status, err := p.service.Verify(r.Context(), paramTxRefValue)
		if err != nil {
			util.ServerErrorResponse(w, err)
			return
		}
		util.OperationSuccessResponse(w, util.Envelope{"verified_status": status})
	}
}

func (p *Payment) ConfirmPaymentHandler(w http.ResponseWriter, r *http.Request) {
	const ParamTxRef = "tx_ref"
	paramValues := r.URL.Query()
	paramTxRefValue := paramValues.Get(ParamTxRef)
	if paramTxRefValue != "" {
		err := p.service.Confirm(r.Context(), paramTxRefValue)
		if err != nil {
			util.ServerErrorResponse(w, err)
			return
		}
		util.OperationSuccessMessageResponse(w, "payment confirmed successfully")
	}
}

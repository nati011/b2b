package application

import (
	"errors"
	"net/http"

	"b2b.nati011.github.com/internal/adapter/primary/rest/handler"
	util "b2b.nati011.github.com/internal/adapter/primary/rest/handler/util"
	application_core "b2b.nati011.github.com/internal/core/application"
	"b2b.nati011.github.com/internal/core/application/payment"
	domain_core "b2b.nati011.github.com/internal/core/domain"
)

var (
	ErrUnknownPaymentCommand = errors.New("unknown user command")
)

type Payment struct {
	service payment.Provider
}

func InitPayment() {
	handler.Register(new(Payment))
}

func (p *Payment) Init(applicationServices *application_core.Container, domainService *domain_core.Container) error {
	p.service = applicationServices.PaymentService
	return nil
}

func (p *Payment) Routes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/payment/webhook/{gateway_id}/{tx_ref}", p.CallbackHandler)

}

func (p *Payment) CallbackHandler(w http.ResponseWriter, r *http.Request) {
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

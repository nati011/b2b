package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"strconv"

	"b2b.nati011.github.com/internal/adapter/primary/rest/handler"
	util "b2b.nati011.github.com/internal/adapter/primary/rest/handler/util"
	application_core "b2b.nati011.github.com/internal/core/application"
	"b2b.nati011.github.com/internal/core/application/payment_partner"
	domain_core "b2b.nati011.github.com/internal/core/domain"
)

type CreatePaymentPartnerRequest struct {
	Name             string `json:"name"`
	Icon             string `json:"icon"`
	Status           string `json:"status"`
	Init_payment_url string `json:"init_payment_url"`
}

type PaymentPartner struct {
	service payment_partner.Provider
}

func InitPaymentPartner() {
	handler.Register(new(PaymentPartner))
}

func (r *PaymentPartner) Init(applicationServices *application_core.Container, domainService *domain_core.Container) error {
	r.service = applicationServices.PaymentPartnerService
	return nil
}

func (p *PaymentPartner) Routes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/payment_option", p.GetPaymentPartnersHandler)
	mux.HandleFunc("GET /api/v1/payment_option/active", p.GetActivePaymentPartnersHandler)
	mux.HandleFunc("POST /api/v1/payment_option", p.CreatePaymentPartnerHandler)
	mux.HandleFunc("PUT /api/v1/payment_option/activate", p.ActivatePaymentPartnerHandler)
	mux.HandleFunc("PUT /api/v1/payment_option/deactivate", p.DectivatePaymentPartnerHandler)
}

func (p *PaymentPartner) GetPaymentPartnersHandler(w http.ResponseWriter, r *http.Request) {
	const ParamId = "id"
	const ParamName = "name"
	const ParamStatus = "status"

	paramValues := r.URL.Query()
	paramNameValue := paramValues.Get(ParamName)
	paramStatusValue := paramValues.Get(ParamStatus)

	paramIdValue := paramValues.Get(ParamId)
	if paramIdValue != "" {
		typedParamId, err := strconv.Atoi(paramIdValue)
		if err != nil {
			util.RequestErrorResponse(w, err)
			return

		}
		resp, err := p.service.Get(r.Context(), typedParamId)
		if err != nil {
			switch err {
			case payment_partner.ErrIdNotFound:
			default:
				util.ServerErrorResponse(w, err)
				return
			}
		}
		util.OperationSuccessResponse(w, util.Envelope{"payment_option": resp})
	} else if paramNameValue != "" || paramStatusValue != "" {
		params := &payment_partner.GetByParamRequest{
			Name:   paramNameValue,
			Status: paramStatusValue,
		}

		resp, err := p.service.GetByParam(r.Context(), params)
		if err != nil {
			switch err {
			case payment_partner.ErrUnknown:
				util.ServerErrorResponse(w, err)
				return
			default:
				util.RequestErrorResponse(w, err)
				return
			}
		}
		util.OperationSuccessResponse(w, util.Envelope{"payment_options": resp})

	} else {
		payment_options, err := p.service.GetAll(r.Context())
		if err != nil {
			switch err {
			default:
				util.ServerErrorResponse(w, err)
				return
			}
		}
		util.OperationSuccessResponse(w, util.Envelope{"payment_options": payment_options})
	}
}

func (p *PaymentPartner) GetActivePaymentPartnersHandler(w http.ResponseWriter, r *http.Request) {
	payment_options, err := p.service.GetActive(r.Context())
	if err != nil {
		switch err {
		default:
			util.ServerErrorResponse(w, err)
			return
		}
	}
	util.OperationSuccessResponse(w, util.Envelope{"payment_options": payment_options})
}

func (p *PaymentPartner) CreatePaymentPartnerHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		util.RequestErrorResponse(w, err)
		return
	}
	defer r.Body.Close()

	var requestBody CreatePaymentPartnerRequest
	if err := json.Unmarshal(body, &requestBody); err != nil {
		util.RequestErrorResponse(w, err)
		return
	}
	id, err := p.service.Create(r.Context(), (*payment_partner.CreateRequest)(&requestBody))
	if err != nil {
		switch err {
		case payment_partner.ErrUnknown:
			util.ServerErrorResponse(w, err)
			return
		default:
			util.RequestErrorResponse(w, err)
			return
		}
	}
	util.OperationSuccessResponse(w, util.Envelope{"product": id})
}

func (p *PaymentPartner) ActivatePaymentPartnerHandler(w http.ResponseWriter, r *http.Request) {
	const ParamId = "id"
	paramValues := r.URL.Query()
	paramIdValue := paramValues.Get(ParamId)

	if paramIdValue != "" {
		typedParamId, err := strconv.Atoi(paramIdValue)
		if err != nil {
			util.RequestErrorResponse(w, err)
			return

		}
		err = p.service.Activate(r.Context(), typedParamId)
		if err != nil {
			switch err {
			case payment_partner.ErrUnknown:
				util.ServerErrorResponse(w, err)
				return
			default:
				util.RequestErrorResponse(w, err)
				return
			}
		}
		util.OperationSuccessResponse(w, util.Envelope{"detail": "Payment Option Activated Successfully"})
	} else {
		util.RequestErrorResponse(w, util.ErrIdRequired)
	}

}

func (p *PaymentPartner) DectivatePaymentPartnerHandler(w http.ResponseWriter, r *http.Request) {
	const ParamId = "id"
	paramValues := r.URL.Query()
	paramIdValue := paramValues.Get(ParamId)

	if paramIdValue != "" {
		typedParamId, err := strconv.Atoi(paramIdValue)
		if err != nil {
			util.RequestErrorResponse(w, err)
			return

		}
		err = p.service.Deactivate(r.Context(), typedParamId)
		if err != nil {
			switch err {
			case payment_partner.ErrUnknown:
				util.ServerErrorResponse(w, err)
				return
			default:
				util.RequestErrorResponse(w, err)
				return
			}
		}
		util.OperationSuccessResponse(w, util.Envelope{"detail": "Payment Option Deactivated Successfully"})
	} else {
		util.RequestErrorResponse(w, util.ErrIdRequired)
	}

}

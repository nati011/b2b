package application

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"strconv"

	"b2b.nati011.github.com/internal/adapter/primary/rest/handler"
	"b2b.nati011.github.com/internal/adapter/primary/rest/handler/middleware"
	util "b2b.nati011.github.com/internal/adapter/primary/rest/handler/util"
	application_core "b2b.nati011.github.com/internal/core/application"
	"b2b.nati011.github.com/internal/core/application/payment_partner"
	domain_core "b2b.nati011.github.com/internal/core/domain"
)

var (
	ErrUnknownPaymentOptionCommand = errors.New("unknown user command")
)

type CreatePaymentPartnerRequest struct {
	Name    string `json:"name"`
	Icon    string `json:"icon"`
	Status  string `json:"status"`
	BaseURL string `json:"base_url"`
	Secret  string `json:"secret"`
}

type GetPaymentPartnerResponse struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Icon    string `json:"icon"`
	Status  string `json:"status"`
	BaseURL string `json:"base_url"`
}

type GetAllPaymentPartnerResponse struct {
	List []GetPaymentPartnerResponse `json:"payment_options"`
}

type GetPaymentOptionByParamRequest struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

type UpdatePaymentOptionRequest struct {
	BaseURL string `json:"base_url"`
	Secret  string `json:"secret"`
}

type PaymentPartner struct {
	authMiddleware middleware.Auth
	service        payment_partner.Provider
}

func InitPaymentPartner() {
	handler.Register(new(PaymentPartner))
}

func (r *PaymentPartner) Init(authMiddleWare *middleware.Auth, applicationServices *application_core.Container, domainService *domain_core.Container) error {
	r.service = applicationServices.PaymentPartnerService
	r.authMiddleware = *applicationServices.AuthMiddleware
	return nil
}

func (p *PaymentPartner) Routes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/payment_option", func(w http.ResponseWriter, r *http.Request) {
		p.authMiddleware.RequireAuthentication(http.HandlerFunc(p.GetPaymentPartnersHandler)).ServeHTTP(w, r)
	})

	mux.HandleFunc("GET /api/v1/payment_option/active", func(w http.ResponseWriter, r *http.Request) {
		p.authMiddleware.RequireAuthentication(http.HandlerFunc(p.GetActivePaymentPartnersHandler)).ServeHTTP(w, r)
	})

	mux.HandleFunc("POST /api/v1/payment_option", func(w http.ResponseWriter, r *http.Request) {
		p.authMiddleware.RequireAuthentication(http.HandlerFunc(p.CreatePaymentPartnerHandler)).ServeHTTP(w, r)
	})

	mux.HandleFunc("PATCH /api/v1/payment_option/{id}/status", func(w http.ResponseWriter, r *http.Request) {
		p.authMiddleware.RequireAuthentication(http.HandlerFunc(p.StatusCommandHandler)).ServeHTTP(w, r)
	})

	mux.HandleFunc("PATCH /api/v1/payment_option/{id}/secret", func(w http.ResponseWriter, r *http.Request) {
		p.authMiddleware.RequireAuthentication(http.HandlerFunc(p.UpdateSecretHandler)).ServeHTTP(w, r)
	})
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
				util.RequestErrorResponse(w, err)
				return
			default:
				util.ServerErrorResponse(w, err)
				return
			}
		}
		util.OperationSuccessResponse(w, util.Envelope{"payment_option": GetPaymentPartnerResponse(resp)})
	} else if paramNameValue != "" || paramStatusValue != "" {
		var response GetAllPaymentPartnerResponse
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
			case payment_partner.ErrEmptyGetContent:
				util.OperationSuccessResponse(w, response)
				return
			default:
				util.RequestErrorResponse(w, err)
				return
			}
		}

		for _, i := range resp.List {
			response.List = append(response.List, GetPaymentPartnerResponse(i))
		}
		util.OperationSuccessResponse(w, response)
	} else {
		resp, err := p.service.GetAll(r.Context())
		if err != nil {
			switch err {
			case payment_partner.ErrEmptyGetContent:
			default:
				util.ServerErrorResponse(w, err)
				return
			}
		}
		var response GetAllPaymentPartnerResponse
		for _, i := range resp.List {
			response.List = append(response.List, GetPaymentPartnerResponse(i))
		}
		util.OperationSuccessResponse(w, response)
	}
}

func (p *PaymentPartner) GetActivePaymentPartnersHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := p.service.GetActive(r.Context())
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
	var response GetAllPaymentPartnerResponse
	for _, i := range resp.List {
		response.List = append(response.List, GetPaymentPartnerResponse(i))
	}
	util.OperationSuccessResponse(w, util.Envelope{"payment_options": response})
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
	util.OperationSuccessResponse(w, util.Envelope{"payment_option": id})
}

const (
	ACTIVATE_PAYMENT_OPTION_COMMAND   = "activate"
	DEACTIVATE_PAYMENT_OPTION_COMMAND = "deactivate"
)

func (p *PaymentPartner) StatusCommandHandler(w http.ResponseWriter, r *http.Request) {
	const ParamCommand = "command"
	paramValues := r.URL.Query()
	paramCommandValue := paramValues.Get(ParamCommand)
	if paramCommandValue != "" {
		typedParamId, err := util.GetPathParam(r, 4)
		if err != nil {
			util.RequestErrorResponse(w, err)
			return
		}
		switch paramCommandValue {
		case ACTIVATE_PAYMENT_OPTION_COMMAND:
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
			util.OperationSuccessMessageResponse(w, "payment option activated successfully")
		case DEACTIVATE_PAYMENT_OPTION_COMMAND:
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
			util.OperationSuccessMessageResponse(w, "payment option deactivated successfully")
		default:
			util.RequestErrorResponse(w, ErrUnknownPaymentOptionCommand)
		}
	}
}

func (p *PaymentPartner) UpdateSecretHandler(w http.ResponseWriter, r *http.Request) {
	typedParamId, err := util.GetPathParam(r, 4)
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
	var requestBody UpdatePaymentOptionRequest
	if err := json.Unmarshal(body, &requestBody); err != nil {
		util.RequestErrorResponse(w, err)
		return
	}
	err = p.service.UpdatePartnerSecret(r.Context(), payment_partner.UpdatePartnerSecret{
		Id:      typedParamId,
		BaseURL: requestBody.BaseURL,
		Secret:  requestBody.Secret,
	})
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
	util.OperationSuccessMessageResponse(w, "secret updated successfully")
}

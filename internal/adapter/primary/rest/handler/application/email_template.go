package application

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"b2b.nati011.github.com/internal/adapter/primary/rest/handler"
	util "b2b.nati011.github.com/internal/adapter/primary/rest/handler/util"
	application_core "b2b.nati011.github.com/internal/core/application"
	"b2b.nati011.github.com/internal/core/application/middleware"
	"b2b.nati011.github.com/internal/core/application/template"
	domain_core "b2b.nati011.github.com/internal/core/domain"
)

type CreateEmailTemplateRequest struct {
	Name         string `json:"name"`
	HtmlTemplate string `json:"temp"`
}

type EmailTemplate struct {
	authMiddleware middleware.Auth
	Service        template.Provider
}

func InitEmailTemplate() {
	handler.Register(new(EmailTemplate))
}

func (e *EmailTemplate) Init(authMiddleWare *middleware.Auth, applicationServices *application_core.Container, domainService *domain_core.Container) error {
	e.authMiddleware = *applicationServices.AuthMiddleware
	e.Service = applicationServices.TemplateService
	return nil
}

func (e *EmailTemplate) Routes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/email_template", func(w http.ResponseWriter, r *http.Request) {
		e.authMiddleware.RequireAuthentication(http.HandlerFunc(e.CreateEmailTemplateHandler)).ServeHTTP(w, r)
	})
	mux.HandleFunc("GET /api/v1/email_template", func(w http.ResponseWriter, r *http.Request) {
		e.authMiddleware.RequireAuthentication(http.HandlerFunc(e.GetEmailTemplateHandler)).ServeHTTP(w, r)
	})
}

func (e *EmailTemplate) CreateEmailTemplateHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		util.RequestErrorResponse(w, err)
		return
	}
	defer r.Body.Close()
	var requestBody CreateEmailTemplateRequest
	if err := json.Unmarshal(body, &requestBody); err != nil {
		util.RequestErrorResponse(w, err)
		return
	}

	id, err := e.Service.Create(r.Context(), &template.CreateRequest{
		Name:         requestBody.Name,
		HtmlTemplate: requestBody.HtmlTemplate,
	})
	if err != nil {
		switch err {
		case template.ErrUnknown:
			util.ServerErrorResponse(w, err)
			return
		default:
			util.RequestErrorResponse(w, err)
			return
		}
	}
	util.OperationSuccessResponse(w, util.Envelope{"email_template": id})
}

func (e *EmailTemplate) GetEmailTemplateHandler(w http.ResponseWriter, r *http.Request) {
	const ParamId = "id"
	paramValues := r.URL.Query()
	paramIdValue := paramValues.Get(ParamId)
	if paramIdValue != "" {
		typedParamId, err := strconv.Atoi(paramIdValue)
		if err != nil {
			util.RequestErrorResponse(w, err)
			return
		}
		resp, err := e.Service.Get(r.Context(), typedParamId)
		if err != nil {
			switch err {
			case template.ErrIdNotFound:
				util.OperationSuccessResponse(w, util.Envelope{"email_template": nil})
				return
			default:
				util.ServerErrorResponse(w, err)
				return
			}
		}
		util.OperationSuccessResponse(w, util.Envelope{"email_template": resp})
	} else {
		resp, err := e.Service.GetAll(r.Context())
		if err != nil {
			switch err {
			case template.ErrEmptyGetContent:
				util.OperationSuccessResponse(w, util.Envelope{"email_template": nil})
			default:
				util.ServerErrorResponse(w, err)
				return
			}
		}
		util.OperationSuccessResponse(w, util.Envelope{"email_template": resp})
	}
}

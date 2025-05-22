package application

import (
	"encoding/json"
	"io"
	"net/http"

	"b2b.nati011.github.com/internal/adapter/primary/rest/handler"
	"b2b.nati011.github.com/internal/adapter/primary/rest/handler/middleware"
	util "b2b.nati011.github.com/internal/adapter/primary/rest/handler/util"
	application_core "b2b.nati011.github.com/internal/core/application"
	"b2b.nati011.github.com/internal/core/application/email"
	"b2b.nati011.github.com/internal/core/application/template"
	domain_core "b2b.nati011.github.com/internal/core/domain"
)

type SendEmailRequest struct {
	To         string            `json:"to"`
	Subject    string            `json:"subject"`
	TemplateId int               `json:"template_id"`
	Args       map[string]string `json:"args"`
	ExternalId string            `json:"ext_id"`
}

type Email struct {
	authMiddleware middleware.Auth
	Service        email.Provider
}

func InitEmail() {
	handler.Register(new(Email))
}

func (e *Email) Init(authMiddleWare *middleware.Auth, applicationServices *application_core.Container, domainService *domain_core.Container) error {
	e.authMiddleware = *applicationServices.AuthMiddleware
	e.Service = applicationServices.EmailService
	return nil
}

func (e *Email) Routes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/email", func(w http.ResponseWriter, r *http.Request) {
		e.authMiddleware.RequireAuthentication(http.HandlerFunc(e.SendEmailHandler)).ServeHTTP(w, r)
	})
}

func (e *Email) SendEmailHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		util.RequestErrorResponse(w, err)
		return
	}
	defer r.Body.Close()
	var requestBody SendEmailRequest
	if err := json.Unmarshal(body, &requestBody); err != nil {
		util.RequestErrorResponse(w, err)
		return
	}

	err = e.Service.Send(&email.SendRequest{
		To:         requestBody.To,
		Subject:    requestBody.Subject,
		TemplateId: requestBody.TemplateId,
		Args:       requestBody.Args,
		ExternalId: requestBody.ExternalId})
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
	util.OperationSuccessMessageResponse(w, "email sent successfuly!")
}

package domain

import (
	"fmt"
	"net/http"
	"strconv"

	"b2b.nati011.github.com/internal/adapter/primary/rest/handler"
	util "b2b.nati011.github.com/internal/adapter/primary/rest/handler/util"
	application_core "b2b.nati011.github.com/internal/core/application"
	"b2b.nati011.github.com/internal/core/application/middleware"
	"b2b.nati011.github.com/internal/core/application/resource"
	domain_core "b2b.nati011.github.com/internal/core/domain"
	"b2b.nati011.github.com/internal/core/domain/invoice"
)

type Invoice struct {
	authMiddleware middleware.Auth
	Service        invoice.Provider
}

func InitInvoice() {
	handler.Register(new(Invoice))

	handler.RegisterResource(resource.CreateRequest{
		Name:     "invoice",
		Action:   "ALL",
		Resource: "/api/v1/invoice",
	})
	handler.RegisterResource(resource.CreateRequest{
		Name:     "invoice_document",
		Action:   "ALL",
		Resource: "/api/v1/invoice/html",
	})
}

func (i *Invoice) Init(authMiddleWare *middleware.Auth, applicationServices *application_core.Container, domainService *domain_core.Container) error {
	i.Service = domainService.InvoiceService
	i.authMiddleware = *applicationServices.AuthMiddleware
	return nil
}

func (i *Invoice) Routes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/invoice", func(w http.ResponseWriter, r *http.Request) {
		i.authMiddleware.RequireAuthentication(http.HandlerFunc(i.GetHandler)).ServeHTTP(w, r)
	})

	mux.HandleFunc("GET /api/v1/invoice/html", func(w http.ResponseWriter, r *http.Request) {
		i.authMiddleware.RequireAuthentication(http.HandlerFunc(i.GetHtmlHandlerHtml)).ServeHTTP(w, r)
	})
}

func (i *Invoice) GetHandler(w http.ResponseWriter, r *http.Request) {
	const ParamOrderId = "order_id"
	paramValues := r.URL.Query()
	paramOrderIdValue := paramValues.Get(ParamOrderId)
	if paramOrderIdValue != "" {
		var typedOrderId int
		var err error
		if paramOrderIdValue != "" {
			typedOrderId, err = strconv.Atoi(paramOrderIdValue)
			if err != nil {
				util.RequestErrorResponse(w, err)
				return
			}
		}
		resp, err := i.Service.GetByParam(r.Context(), &invoice.GetByParamRequest{
			OrderId: typedOrderId,
		})
		if err != nil {
			switch err {
			case invoice.ErrSysUnknown:
				util.ServerErrorResponse(w, err)
				return
			case invoice.ErrSysEmptyGetContent:
			default:
				util.RequestErrorResponse(w, err)
				return
			}
		}
		util.OperationSuccessResponse(w, util.Envelope{"invoice": resp})
	}
}

func (i *Invoice) GetHtmlHandlerHtml(w http.ResponseWriter, r *http.Request) {
	const ParamOrderId = "order_id"
	paramValues := r.URL.Query()
	paramOrderIdValue := paramValues.Get(ParamOrderId)
	if paramOrderIdValue != "" {
		var typedOrderId int
		var err error
		if paramOrderIdValue != "" {
			typedOrderId, err = strconv.Atoi(paramOrderIdValue)
			if err != nil {
				util.RequestErrorResponse(w, err)
				return
			}
		}
		resp, err := i.Service.GetInvoiceHtml(r.Context(), typedOrderId)

		if err != nil {
			switch err {
			case invoice.ErrSysUnknown:
				util.ServerErrorResponse(w, err)
				return
			case invoice.ErrSysEmptyGetContent:
			default:
				util.RequestErrorResponse(w, err)
				return
			}
		}
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, resp.HTML)
		return
	}
}

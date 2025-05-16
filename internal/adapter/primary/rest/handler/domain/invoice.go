package domain

import (
	"net/http"
	"strconv"

	"b2b.nati011.github.com/internal/adapter/primary/rest/handler"
	util "b2b.nati011.github.com/internal/adapter/primary/rest/handler/util"
	application_core "b2b.nati011.github.com/internal/core/application"
	domain_core "b2b.nati011.github.com/internal/core/domain"
	"b2b.nati011.github.com/internal/core/domain/invoice"
)

type Invoice struct {
	Service invoice.Provider
}

func InitInvoice() {
	handler.Register(new(Invoice))
}

func (i *Invoice) Init(applicationServices *application_core.Container, domainService *domain_core.Container) error {
	i.Service = domainService.InvoiceService
	return nil
}

func (i *Invoice) Routes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/invoice", i.GetHandler)
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

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
	service invoice.Provider
}

func InitInvoice() {
	handler.Register(new(Invoice))
}

func (i *Invoice) Init(applicationServices *application_core.Container, domainServices *domain_core.Container) error {
	i.service = domainServices.InvoiceService
	return nil
}

func (i *Invoice) Routes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/invoice", i.GetHandler)
}

func (i *Invoice) GetHandler(w http.ResponseWriter, r *http.Request) {
	const ParamId = "id"
	const ParamStatus = "status"
	const ParamExternalId = "external_id"
	const ParamOrderId = "order_id"
	const ParamCreated_Date = "created_date"

	paramValues := r.URL.Query()
	paramStatusValue := paramValues.Get(ParamStatus)
	ParamExternalIdValue := paramValues.Get(ParamExternalId)
	ParamOrderIdValue := paramValues.Get(ParamOrderId)

	paramIdValue := paramValues.Get(ParamId)
	if paramIdValue != "" {
		typedParamId, err := strconv.Atoi(paramIdValue)
		if err != nil {
			util.RequestErrorResponse(w, err)
			return

		}
		resp, err := i.service.Get(r.Context(), typedParamId)
		if err != nil {
			switch err {
			case invoice.ErrSysIdNotFound:
			default:
				util.ServerErrorResponse(w, err)
				return
			}
		}
		util.OperationSuccessResponse(w, util.Envelope{"product": resp})

	} else if ParamCategoryIdValue != "" || ParamPriceMinValue != "" || ParamPriceMaxValue != "" {

	} else {
		resp, err := i.service.GetAll(r.Context())
		if err != nil {
			switch err {
			case invoice.ErrSysEmptyGetContent:
				util.RequestErrorResponse(w, err)
				return
			default:
				util.ServerErrorResponse(w, err)
				return
			}
		}
		util.OperationSuccessResponse(w, util.Envelope{"invoices": resp})
	}
}

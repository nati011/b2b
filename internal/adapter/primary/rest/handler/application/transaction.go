package application

import (
	"net/http"
	"time"

	"strconv"

	"b2b.nati011.github.com/internal/adapter/primary/rest/handler"
	util "b2b.nati011.github.com/internal/adapter/primary/rest/handler/util"
	application_core "b2b.nati011.github.com/internal/core/application"
	"b2b.nati011.github.com/internal/core/application/transaction"
	domain_core "b2b.nati011.github.com/internal/core/domain"
)

type GetByParamRequest struct {
	Date      time.Time `json:"date"`
	PartnerId int       `json:"partner_id"`
	TxRef     string    `json:"tx_ref"`
	Status    string    `json:"status"`
}

type GetResponse struct {
	Id        int       `json:"id"`
	Date      time.Time `json:"date"`
	Amount    float64   `json:"amount"`
	PartnerId int       `json:"partner_id"`
	TxRef     string    `json:"tx_ref"`
	Status    string    `json:"status"`
}

type GetAllResponse struct {
	List []GetResponse `json:"transactions"`
}

type Transaction struct {
	service transaction.Provider
}

func InitTransaction() {
	handler.Register(new(Transaction))
}

func (r *Transaction) Init(applicationServices *application_core.Container, domainService *domain_core.Container) error {
	r.service = applicationServices.TransactionService
	return nil
}

func (p *Transaction) Routes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/transaction", p.GetTransactionsHandler)
}

func (p *Transaction) GetTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	const ParamId = "id"
	const ParamPartnerId = "partner_id"
	const ParamDate = "date"

	paramValues := r.URL.Query()
	paramPartnerIdValue := paramValues.Get(ParamPartnerId)
	paramDateValue := paramValues.Get(ParamDate)

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
			case transaction.ErrIdNotFound:
			default:
				util.ServerErrorResponse(w, err)
				return
			}
		}
		util.OperationSuccessResponse(w, util.Envelope{"transaction": GetResponse(resp)})
	} else if paramPartnerIdValue != "" || paramDateValue != "" {
		var response GetAllResponse
		typedPartnerId, err := strconv.Atoi(paramPartnerIdValue)
		if err != nil {
			util.RequestErrorResponse(w, err)
			return
		}

		parsedDate, err := time.Parse(paramDateValue, "2024-09-19 14:00:00")
		if err != nil {
			util.RequestErrorResponse(w, err)
			return
		}

		params := &transaction.GetByParamRequest{
			PartnerId: typedPartnerId,
			Date:      parsedDate,
		}

		resp, err := p.service.GetByParam(r.Context(), params)
		if err != nil {
			switch err {
			case transaction.ErrEmptyGetContent:
				util.OperationSuccessResponse(w, response)
				return
			case transaction.ErrUnknown:
				util.ServerErrorResponse(w, err)
				return
			default:
				util.RequestErrorResponse(w, err)
				return
			}
		}
		for _, i := range resp.List {
			response.List = append(response.List, GetResponse(i))
		}
		util.OperationSuccessResponse(w, response)

	} else {
		var response GetAllResponse
		resp, err := p.service.GetAll(r.Context())
		if err != nil {
			switch err {
			case transaction.ErrEmptyGetContent:
				util.OperationSuccessResponse(w, response)
				return
			case transaction.ErrUnknown:
				util.ServerErrorResponse(w, err)
				return
			default:
				util.RequestErrorResponse(w, err)
				return
			}
		}
		for _, i := range resp.List {
			response.List = append(response.List, GetResponse(i))
		}
		util.OperationSuccessResponse(w, response)
	}
}

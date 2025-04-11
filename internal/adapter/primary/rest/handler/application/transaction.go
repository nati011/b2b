package handler

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

type CreateTransactionRequest struct {
	User_Id    int   `json:"user_id"`
	Amount     int64 `json:"amount"`
	Partner_Id int   `json:"partner_id"`
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
	const ParamUserId = "user_id"
	const ParamDate = "date"

	paramValues := r.URL.Query()
	paramPartnerIdValue := paramValues.Get(ParamPartnerId)
	paramDateValue := paramValues.Get(ParamDate)
	paramUserIdValue := paramValues.Get(ParamUserId)

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
		util.OperationSuccessResponse(w, util.Envelope{"transaction": resp})
	} else if paramPartnerIdValue != "" || paramDateValue != "" || paramUserIdValue != "" {
		typedPartnerId, err := strconv.Atoi(paramPartnerIdValue)
		if err != nil {
			util.RequestErrorResponse(w, err)
			return
		}
		typedUserId, err := strconv.Atoi(paramUserIdValue)
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
			Partner_Id: typedPartnerId,
			Date:       parsedDate,
			User_Id:    typedUserId,
		}

		resp, err := p.service.GetByParam(r.Context(), params)
		if err != nil {
			switch err {
			case transaction.ErrUnknown:
				util.ServerErrorResponse(w, err)
				return
			default:
				util.RequestErrorResponse(w, err)
				return
			}
		}
		util.OperationSuccessResponse(w, util.Envelope{"transactions": resp})

	} else {
		transactions, err := p.service.GetAll(r.Context())
		if err != nil {
			switch err {
			default:
				util.ServerErrorResponse(w, err)
				return
			}
		}
		util.OperationSuccessResponse(w, util.Envelope{"transactions": transactions})
	}
}

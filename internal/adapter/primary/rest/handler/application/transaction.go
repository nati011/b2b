package handler

import (
	"encoding/json"
	"io"
	"log"
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
	mux.HandleFunc("POST /api/v1/transaction", p.CreateTransactionHandler)
	// mux.HandleFunc("PUT /api/v1/transaction/activate", p.ActivateTransactionHandler)
	// mux.HandleFunc("PUT /api/v1/transaction/deactivate", p.DectivateTransactionHandler)
}

func (p *Transaction) GetTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	const ParamId = "id"
	const ParamPartnerId = "partner_id"
	const ParamDate = "date"
	const ParamUserId = "user_id"

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
		var typedPartnerId int
		var typedUserId int
		var parsedDate time.Time

		if paramPartnerIdValue != "" {
			typedPartnerId, _ = strconv.Atoi(paramPartnerIdValue)
			// if err != nil {
			// 	util.RequestErrorResponse(w, err)
			// 	return

			// }
		}

		if paramUserIdValue != "" {
			typedUserId, _ = strconv.Atoi(paramUserIdValue)
			// if err != nil {
			// 	util.RequestErrorResponse(w, err)
			// 	return

			// }
		}

		if paramDateValue != "" {
			parsedDate, err := time.Parse(paramDateValue, "2024-09-19 14:00:00")
			log.Print(parsedDate.String())
			if err != nil {
				util.RequestErrorResponse(w, err)
				return

			}
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

func (t *Transaction) CreateTransactionHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		util.RequestErrorResponse(w, err)
		return
	}
	defer r.Body.Close()

	var requestBody CreateTransactionRequest
	if err := json.Unmarshal(body, &requestBody); err != nil {
		util.RequestErrorResponse(w, err)
		return
	}
	id, err := t.service.Create(r.Context(), (*transaction.CreateRequest)(&requestBody))
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
	util.OperationSuccessResponse(w, util.Envelope{"transaction": id})
}

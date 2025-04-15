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

type CreateTransactionRequest struct {
	User_Id    int   `json:"user_id"`
	Amount     int64 `json:"amount"`
	Partner_Id int   `json:"partner_id"`
}

type GetByParamRequest struct {
	Date       time.Time `json:"date"`
	Partner_Id int       `json:"partner_id"`
	User_Id    int       `json:"user_id"`
}

type GetResponse struct {
	Id         int       `json:"id"`
	User_Id    int       `json:"user_id"`
	Date       time.Time `json:"date"`
	Amount     int64     `json:"amount"`
	Partner_Id int       `json:"partner_id"`
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
	const ParamUserId = "user_id"
	const ParamLimit = "limit"
	const ParamOffset = "offset"
	const ParamDate = "date"

	paramValues := r.URL.Query()
	paramPartnerIdValue := paramValues.Get(ParamPartnerId)
	paramDateValue := paramValues.Get(ParamDate)
	paramUserIdValue := paramValues.Get(ParamUserId)
	paramLimitValue := paramValues.Get(ParamLimit)
	paramOffsetValue := paramValues.Get(ParamOffset)

	var typedLimit int
	var typedOffset int
	var err error

	if paramLimitValue != "" {
		typedLimit, err = strconv.Atoi(paramLimitValue)
		if err != nil {
			util.RequestErrorResponse(w, err)
		}
	}
	if paramOffsetValue != "" {
		typedOffset, err = strconv.Atoi(paramOffsetValue)

		if err != nil {
			util.RequestErrorResponse(w, err)
		}
	}

	pagination := transaction.Pagination{
		Limit:  typedLimit,
		Offset: typedOffset,
	}

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
	} else if paramPartnerIdValue != "" || paramDateValue != "" || paramUserIdValue != "" {
		var typedPartnerId int
		var typedUserId int
		var parsedDate time.Time
		var err error

		if paramPartnerIdValue != "" {
			typedPartnerId, err = strconv.Atoi(paramPartnerIdValue)
			if err != nil {
				util.RequestErrorResponse(w, err)
				return

			}
		}

		if paramUserIdValue != "" {
			typedUserId, err = strconv.Atoi(paramUserIdValue)
			if err != nil {
				util.RequestErrorResponse(w, err)
				return

			}

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

		resp, err := p.service.GetByParam(r.Context(), params, &pagination)
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
		var response GetAllResponse
		for _, i := range resp.List {
			response.List = append(response.List, GetResponse(i))
		}
		util.OperationSuccessResponse(w, util.Envelope{"transactions": response})

	} else {
      resp,err:=p.service.GetAll(r.Context(), &pagination)
		if err != nil {
			switch err {
			default:
				util.ServerErrorResponse(w, err)
				return
			}
		}
		var response GetAllResponse
		for _, i := range resp.List {
			response.List = append(response.List, GetResponse(i))
		}
		util.OperationSuccessResponse(w, util.Envelope{"transactions": response})
	}
}

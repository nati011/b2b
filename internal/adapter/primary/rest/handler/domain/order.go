package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"b2b.nati011.github.com/internal/adapter/primary/rest/handler"
	util "b2b.nati011.github.com/internal/adapter/primary/rest/handler/util"
	application_core "b2b.nati011.github.com/internal/core/application"
	domain_core "b2b.nati011.github.com/internal/core/domain"

	"b2b.nati011.github.com/internal/core/domain/order"
	"b2b.nati011.github.com/internal/core/domain/product"
)

type OrderItem struct {
	ProductId int `json:"id"`
	Quantity  int `json:"quantity"`
}

type PlaceOrderRequest struct {
	RetailerId int         `json:"retailer_id"`
	Items      []OrderItem `json:"items"`
}

type GetOrderResponse struct {
	Id         int         `json:"id"`
	RetailerId int         `json:"retailer_id"`
	Items      []OrderItem `json:"items"`
	Total      float32     `json:"total"`
	Status     string      `json:"status"`
}

type GetAllOrderResponse struct {
	List []GetOrderResponse `json:"orders"`
}

type GetOrderByParamRequest struct {
	RetailerId int    `json:"retailer_id"`
	Status     string `json:"status"`
}

var (
	ErrUnknownCommand = errors.New("unknown command")
)

func InitOrder() {
	handler.Register(new(Order))
}

type Order struct {
	service order.Provider
}

func (r *Order) Init(applicationServices *application_core.Container, domainService *domain_core.Container) error {
	r.service = domainService.OrderService
	return nil
}

func (o *Order) Routes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/order", o.GetHandler)
	mux.HandleFunc("POST /api/v1/order", o.PostHandler)
	mux.HandleFunc("PATCH /api/v1/order", o.CommandHandler)
}

func (o *Order) GetHandler(w http.ResponseWriter, r *http.Request) {
	// get by id
	const ParamId = "id"
	const ParamRetailerId = "retailer_id"
	const ParamStatus = "status"

	paramValues := r.URL.Query()
	paramIdValue := paramValues.Get(ParamId)
	paramRetailerIdValue := paramValues.Get(ParamRetailerId)
	paramStatus := paramValues.Get(ParamStatus)
	if paramIdValue != "" {
		typedParamId, err := strconv.Atoi(paramIdValue)
		if err != nil {
			util.RequestErrorResponse(w, err)
			return
		}

		resp, err := o.service.Get(r.Context(), typedParamId)
		if err != nil {
			switch err {
			case order.ErrUnknown:
				util.ServerErrorResponse(w, err)
			default:
				util.RequestErrorResponse(w, err)
				return
			}
		}
		util.OperationSuccessResponse(w, util.Envelope{"order": resp})
	} else if paramRetailerIdValue != "" || paramStatus != "" {
		// get by param
		var typedRetailerId int
		var err error
		if paramRetailerIdValue != "" {
			typedRetailerId, err = strconv.Atoi(paramRetailerIdValue)
			if err != nil {
				util.RequestErrorResponse(w, err)
				return
			}
		}
		resp, err := o.service.GetByParam(r.Context(), &order.GetByParamRequest{
			RetailerId: typedRetailerId,
			Status:     paramStatus,
		})
		if err != nil {
			switch err {
			case order.ErrUnknown:
				util.ServerErrorResponse(w, err)
				return
			default:
				util.RequestErrorResponse(w, err)
				return
			}
		}
		util.OperationSuccessResponse(w, util.Envelope{"orders": resp})
	} else {
		// get all
		resp, err := o.service.GetAll(r.Context())
		if err != nil {
			switch err {
			case order.ErrEmptyGetResponse:
				util.RequestErrorResponse(w, err)
				return
			default:
				util.ServerErrorResponse(w, err)
				return
			}
		}
		util.OperationSuccessResponse(w, util.Envelope{"orders": resp})
	}

}

func (o *Order) PostHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		util.RequestErrorResponse(w, err)
		return
	}
	defer r.Body.Close()

	var requestBody PlaceOrderRequest
	if err := json.Unmarshal(body, &requestBody); err != nil {
		util.RequestErrorResponse(w, err)
		return
	}
	var orderItems []order.Item
	for _, i := range requestBody.Items {
		orderItems = append(orderItems, (order.Item)(i))
	}
	id, err := o.service.Place(r.Context(), &order.PlaceRequest{
		RetailerId: requestBody.RetailerId,
		Items:      orderItems,
	})
	if err != nil {
		switch err {
		case order.ErrUnknown:
			util.ServerErrorResponse(w, err)
			return
		default:
			util.RequestErrorResponse(w, err)
			return
		}
	}
	util.OperationSuccessResponse(w, util.Envelope{"order": id})
}

const (
	CANCEL_COMMAND = "cancel"
)

func (p *Order) CommandHandler(w http.ResponseWriter, r *http.Request) {
	const ParamId = "id"
	const ParamCommand = "command"
	paramValues := r.URL.Query()
	paramCommandValue := paramValues.Get(ParamCommand)
	paramIdValue := paramValues.Get(ParamId)

	if paramCommandValue != "" && paramIdValue != "" {
		typedParamId, err := strconv.Atoi(paramIdValue)
		if err != nil {
			util.RequestErrorResponse(w, err)
			return
		}
		switch paramCommandValue {
		case CANCEL_COMMAND:
			err = p.service.Cancel(r.Context(), typedParamId)
			if err != nil {
				switch err {
				case product.ErrUnknown:
					util.ServerErrorResponse(w, err)
					return
				default:
					util.RequestErrorResponse(w, err)
					return
				}
			}
		default:
			util.RequestErrorResponse(w, ErrUnknownCommand)
		}
		util.OperationSuccessResponse(w, util.Envelope{"order": typedParamId})
	}
}

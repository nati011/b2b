package domain

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"b2b.nati011.github.com/internal/adapter/primary/rest/handler"
	util "b2b.nati011.github.com/internal/adapter/primary/rest/handler/util"
	application_core "b2b.nati011.github.com/internal/core/application"
	"b2b.nati011.github.com/internal/core/application/middleware"
	domain_core "b2b.nati011.github.com/internal/core/domain"

	"b2b.nati011.github.com/internal/core/domain/order"
	"b2b.nati011.github.com/internal/core/domain/product"
)

type OrderItem struct {
	ProductId   int    `json:"id"`
	ProductName string `json:"name"`
	Quantity    int    `json:"quantity"`
}

type PlaceOrderRequest struct {
	RetailerId       int         `json:"retailer_id"`
	Items            []OrderItem `json:"items"`
	PaymentPartnerId int         `json:"payment_partner_id"`
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
	authMiddleware middleware.Auth
	service        order.Provider
}

func (o *Order) Init(authMiddleWare *middleware.Auth, applicationServices *application_core.Container, domainService *domain_core.Container) error {
	o.service = domainService.OrderService
	o.authMiddleware = *applicationServices.AuthMiddleware
	return nil
}

func (o *Order) Routes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/order", func(w http.ResponseWriter, r *http.Request) {
		o.authMiddleware.RequireAuthentication(http.HandlerFunc(o.GetHandler)).ServeHTTP(w, r)
	})

	mux.HandleFunc("POST /api/v1/order", func(w http.ResponseWriter, r *http.Request) {
		o.authMiddleware.RequireAuthentication(http.HandlerFunc(o.PostHandler)).ServeHTTP(w, r)
	})

	mux.HandleFunc("POST /api/v1/order/init_settlement", func(w http.ResponseWriter, r *http.Request) {
		o.authMiddleware.RequireAuthentication(http.HandlerFunc(o.InitPaymentHandler)).ServeHTTP(w, r)
	})

	mux.HandleFunc("PATCH /api/v1/order", func(w http.ResponseWriter, r *http.Request) {
		o.authMiddleware.RequireAuthentication(http.HandlerFunc(o.CommandHandler)).ServeHTTP(w, r)
	})

	mux.HandleFunc("GET /api/v1/orders/retailer", func(w http.ResponseWriter, r *http.Request) {
		o.authMiddleware.RequireAuthentication(http.HandlerFunc(o.GetRetailerOrders)).ServeHTTP(w, r)
	})

	mux.HandleFunc("GET /api/v1/orders/distributor", func(w http.ResponseWriter, r *http.Request) {
		o.authMiddleware.RequireAuthentication(http.HandlerFunc(o.GetDistributorOrders)).ServeHTTP(w, r)
	})
}

func (o *Order) GetRetailerOrders(w http.ResponseWriter, r *http.Request) {
	const ParamRetailerId = "retailer_id"

	paramValues := r.URL.Query()
	paramRetailerIdValue := paramValues.Get(ParamRetailerId)
	if paramRetailerIdValue != "" {
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
		resp, err := o.service.GetRetailerOrders(r.Context(), typedRetailerId)
		if err != nil {
			switch err {
			case order.ErrUnknown:
				util.ServerErrorResponse(w, err)
				return
			case order.ErrEmptyGetResponse:
			default:
				util.RequestErrorResponse(w, err)
				return
			}
		}
		util.OperationSuccessResponse(w, util.Envelope{"orders": resp})
	}
}

func (o *Order) GetDistributorOrders(w http.ResponseWriter, r *http.Request) {
	const ParamDistributorId = "distributor_id"

	paramValues := r.URL.Query()
	paramDistributorIdValue := paramValues.Get(ParamDistributorId)
	if paramDistributorIdValue != "" {
		var typedDistributorId int
		var err error
		if paramDistributorIdValue != "" {
			typedDistributorId, err = strconv.Atoi(paramDistributorIdValue)
			if err != nil {
				util.RequestErrorResponse(w, err)
				return
			}
		}
		resp, err := o.service.GetDistributorOrders(r.Context(), typedDistributorId)
		if err != nil {
			switch err {
			case order.ErrUnknown:
				util.ServerErrorResponse(w, err)
				return
			case order.ErrEmptyGetResponse:
			default:
				util.RequestErrorResponse(w, err)
				return
			}
		}
		util.OperationSuccessResponse(w, util.Envelope{"orders": resp})
	}
}

func (o *Order) GetHandler(w http.ResponseWriter, r *http.Request) {
	const ParamId = "id"
	const ParamStatus = "status"

	paramValues := r.URL.Query()
	paramIdValue := paramValues.Get(ParamId)
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
	} else if paramStatus != "" {
		var err error
		resp, err := o.service.GetByParam(r.Context(), &order.GetByParamRequest{
			Status: paramStatus,
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
		util.OperationSuccessResponse(w, resp)
	} else {
		resp, err := o.service.GetAll(r.Context())
		if err != nil {
			switch err {
			case order.ErrEmptyGetResponse:
			default:
				util.ServerErrorResponse(w, err)
				return
			}
		}
		util.OperationSuccessResponse(w, resp)
	}

}

func (o *Order) InitPaymentHandler(w http.ResponseWriter, r *http.Request) {
	const ParamId = "id"
	paramValues := r.URL.Query()
	paramIdValue := paramValues.Get(ParamId)
	if paramIdValue != "" {
		typedParamId, err := strconv.Atoi(paramIdValue)
		if err != nil {
			util.RequestErrorResponse(w, err)
			return
		}

		id, err := o.service.InitPayment(r.Context(), typedParamId)
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
		util.OperationSuccessResponse(w, util.Envelope{"checkoutUrl": id})
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
		RetailerId:       requestBody.RetailerId,
		Items:            orderItems,
		PaymentPartnerId: requestBody.PaymentPartnerId,
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
			return
		}
		util.OperationSuccessResponse(w, util.Envelope{"order": typedParamId})
	}
}

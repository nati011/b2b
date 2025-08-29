package domain

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"b2b.nati011.github.com/internal/adapter/primary/rest/handler"
	util "b2b.nati011.github.com/internal/adapter/primary/rest/handler/util"
	application_core "b2b.nati011.github.com/internal/core/application"
	"b2b.nati011.github.com/internal/core/application/middleware"
	"b2b.nati011.github.com/internal/core/application/resource"
	domain_core "b2b.nati011.github.com/internal/core/domain"
	distributor_Subscription "b2b.nati011.github.com/internal/core/domain/distributor_subscription"
)

type CreatePlanRequest struct {
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	TermInMonth int     `json:"term_in_month"`
	Description string  `json:"description"`
}

type PlaceRequest struct {
	SubscriptionPlanId int `json:"subscription_plan_id"`
	DistributorId      int `json:"distributor_id"`
	PaymentPartnerId   int `json:"payment_partner_id"`
}

type Distributor_Subscription struct {
	service        distributor_Subscription.Prodvider
	authMiddleware middleware.Auth
}

func InitDistributorSubscription() {
	handler.Register(new(Distributor_Subscription))

	handler.RegisterResource(resource.CreateRequest{
		Name:     "subscription",
		Action:   "ALL",
		Resource: "/api/v1/subscription",
	})

	handler.RegisterResource(resource.CreateRequest{
		Name:     "plan",
		Action:   "ALL",
		Resource: "/api/v1/plan",
	})
}

func (d *Distributor_Subscription) Init(authMiddleWare *middleware.Auth, applicationServices *application_core.Container, domainServices *domain_core.Container) error {
	d.service = domainServices.DistributorSubscriptionService
	d.authMiddleware = *applicationServices.AuthMiddleware
	return nil
}

func (d *Distributor_Subscription) Routes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/subscription", func(w http.ResponseWriter, r *http.Request) {
		d.authMiddleware.RequireAuthentication(http.HandlerFunc(d.GetSubscriptionsHandler)).ServeHTTP(w, r)
	})

	mux.HandleFunc("POST /api/v1/subscription", func(w http.ResponseWriter, r *http.Request) {
		d.authMiddleware.RequireNoAuthentication(http.HandlerFunc(d.CreateSubscriptionHandler)).ServeHTTP(w, r)
	})

	mux.HandleFunc("GET /api/v1/plan", func(w http.ResponseWriter, r *http.Request) {
		d.authMiddleware.RequireNoAuthentication(http.HandlerFunc(d.GetSubscriptionPlanHandler)).ServeHTTP(w, r)
	})

	mux.HandleFunc("POST /api/v1/plan", func(w http.ResponseWriter, r *http.Request) {
		d.authMiddleware.RequireAuthentication(http.HandlerFunc(d.CreateSubscriptionPlanHandler)).ServeHTTP(w, r)
	})
}

func (d *Distributor_Subscription) GetSubscriptionsHandler(w http.ResponseWriter, r *http.Request) {
	const ParamDistributorId = "distributor_id"
	paramValues := r.URL.Query()
	paramDistributorIdValue := paramValues.Get(ParamDistributorId)
	if paramDistributorIdValue != "" {
		typedParamId, err := strconv.Atoi(paramDistributorIdValue)
		if err != nil {
			util.RequestErrorResponse(w, err)
			return
		}
		resp_sub, err := d.service.GetSubscriptionByDistributorId(r.Context(), typedParamId)
		if err != nil {
			switch err {
			case distributor_Subscription.ErrEmptyGetContent:
			default:
				util.ServerErrorResponse(w, err)
				return
			}
		}
		util.OperationSuccessResponse(w, util.Envelope{"sub": resp_sub})
	} else {
		resp_subs, err := d.service.GetSubscriptions(r.Context())
		if err != nil {
			switch err {
			case distributor_Subscription.ErrEmptyGetContent:
			default:
				util.ServerErrorResponse(w, err)
				return
			}
		}
		util.OperationSuccessResponse(w, util.Envelope{"subs": resp_subs})
	}
}

func (d *Distributor_Subscription) GetSubscriptionPlanHandler(w http.ResponseWriter, r *http.Request) {
	const ParamId = "id"
	paramValues := r.URL.Query()
	paramIdValue := paramValues.Get(ParamId)
	if paramIdValue != "" {
		typedParamId, err := strconv.Atoi(paramIdValue)
		if err != nil {
			util.RequestErrorResponse(w, err)
			return
		}
		resp_plan, err := d.service.GetPlan(r.Context(), typedParamId)
		if err != nil {
			switch err {
			case distributor_Subscription.ErrEmptyGetContent:
			default:
				util.ServerErrorResponse(w, err)
				return
			}
		}
		util.OperationSuccessResponse(w, util.Envelope{"plan": resp_plan})
	} else {
		resp_plans, err := d.service.GetAllPlan(r.Context())
		if err != nil {
			switch err {
			case distributor_Subscription.ErrEmptyGetContent:
			default:
				util.ServerErrorResponse(w, err)
				return
			}
		}
		util.OperationSuccessResponse(w, util.Envelope{"plans": resp_plans})
	}
}

func (d *Distributor_Subscription) CreateSubscriptionPlanHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll((r.Body))
	if err != nil {
		util.RequestErrorResponse(w, err)
		return
	}
	defer r.Body.Close()

	var requestBody CreatePlanRequest

	if err := json.Unmarshal(body, &requestBody); err != nil {
		util.RequestErrorResponse(w, err)
		return
	}

	create_req := distributor_Subscription.CreatePlanRequest(requestBody)
	resp_id, err := d.service.CreatePlan(r.Context(), &create_req)
	if err != nil {
		switch err {
		case distributor_Subscription.ErrUnknown:
			util.ServerErrorResponse(w, err)
			return
		default:
			util.RequestErrorResponse(w, err)
			return
		}
	}
	util.OperationSuccessResponse(w, util.Envelope{"id": resp_id})
}

func (d *Distributor_Subscription) CreateSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll((r.Body))
	if err != nil {
		util.RequestErrorResponse(w, err)
		return
	}
	defer r.Body.Close()

	var requestBody PlaceRequest

	if err := json.Unmarshal(body, &requestBody); err != nil {
		util.RequestErrorResponse(w, err)
		return
	}
	create_req := distributor_Subscription.PlaceRequest(requestBody)

	resp_id, err := d.service.Place(r.Context(), &create_req)
	if err != nil {
		switch err {
		case distributor_Subscription.ErrUnknown:
			util.ServerErrorResponse(w, err)
			return
		default:
			util.RequestErrorResponse(w, err)
			return
		}
	}
	util.OperationSuccessResponse(w, util.Envelope{"id": resp_id})
}

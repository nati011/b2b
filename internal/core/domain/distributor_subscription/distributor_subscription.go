package distributor_subscription

import (
	"context"
	"errors"
	"log"

	"b2b.nati011.github.com/internal/core/domain/distributor"
	port_commons "b2b.nati011.github.com/internal/port/commons/db"
	port "b2b.nati011.github.com/internal/port/domain/distributor_subscription"
)

var (
	ErrEmptyGetContent       = errors.New(" empty conent")
	ErrNameMandatory         = errors.New(" name mandatory")
	ErrDescriptionMandatory  = errors.New(" desc mandatory")
	ErrTermCannotBeZero      = errors.New(" term cannot be zero")
	ErrTermCannotBeNegative  = errors.New(" term cannot be negative")
	ErrPriceCannotBeZero     = errors.New(" price cannot be zero")
	ErrPriceCannotBeNegative = errors.New(" price cannot be negative")
	ErrUnknown               = errors.New(" unknown error")
)

type PlaceRequest struct {
	SubscriptionPlanId int
	DistributorId      int
	Status             string
}

type GetSubscriptionResponse struct {
	Id                 int
	SubscriptionPlanId int
	DistributorId      int
	Status             string
}

type GetAllSubscriptionResponse struct {
	List []GetAllSubscriptionPlanResponse
}

type SubscribeResponse struct {
	Id          int    `json:"id"`
	CheckoutUrl string `json:"checkout_url"`
	TxRef       string `json:"tx_ref"`
}

type GetSubscriptionPlanResponse struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	TermInMonth int     `json:"term_in_month"`
	Description string  `json:"desc"`
}

type GetAllSubscriptionPlanResponse struct {
	List []GetSubscriptionPlanResponse `json:"list"`
}

type CreatePlanRequest struct {
	Name        string
	Price       float64
	TermInMonth int
	Description string
}

// TODO: acitvate and deactivate plans
// TODO: make name, price, termInMonth unique
// Note: instead of updating plans, deactivate plan of choice and create a new one instead
type Prodvider interface {
	CreatePlan(ctx context.Context, req *CreatePlanRequest) (int, error)
	GetPlan(ctx context.Context) (GetAllSubscriptionPlanResponse, error)
	Place(ctx context.Context, req *PlaceRequest) (SubscribeResponse, error)
	InitPayment(ctx context.Context, distId int) (SubscribeResponse, error)
	GetSubscriptions(ctx context.Context) (GetAllSubscriptionPlanResponse, error)
}

type DistributorSubscriptionService struct {
	DB port.DB
	DS distributor.Provider
}

func NewDistributorSubscriptionService(
	db port.DB,
	ds distributor.Provider) Prodvider {
	return &DistributorSubscriptionService{
		DB: db,
		DS: ds,
	}
}

func (d *DistributorSubscriptionService) CreatePlan(ctx context.Context, req *CreatePlanRequest) (int, error) {
	if err := d.validateCreate(req); err != nil {
		return 0, err
	}
	id, err := d.DB.CreatePlan(ctx, &port.CreatePlanRequest{
		Name:        req.Name,
		Price:       req.Price,
		TermInMonth: req.TermInMonth,
		Description: req.Description,
	})
	if err != nil {
		return 0, ErrUnknown
	}
	return id, nil
}

func (d *DistributorSubscriptionService) GetPlan(ctx context.Context) (GetAllSubscriptionPlanResponse, error) {
	resp, err := d.DB.GetPlans(ctx)
	if err != nil {
		switch err {
		case port_commons.ErrSysNoRows:
			return GetAllSubscriptionPlanResponse{}, ErrEmptyGetContent
		default:
			log.Printf("failed to get distributor subscription plans err: %v", err.Error())
			return GetAllSubscriptionPlanResponse{}, ErrUnknown
		}
	}
	response := GetAllSubscriptionPlanResponse{}
	for _, i := range resp.List {
		response.List = append(response.List, GetSubscriptionPlanResponse{
			Id:          i.Id,
			Name:        i.Name,
			Price:       i.Price,
			TermInMonth: i.TermInMonth,
			Description: i.Description,
		})
	}
	return response, nil
}

func (d *DistributorSubscriptionService) Place(ctx context.Context, req *PlaceRequest) (SubscribeResponse, error) {
	//validate dist
	//validate subscription plan
	//validate if subscription already exists
	//create subscription record
	//init subscription checkout and return
	return SubscribeResponse{}, nil
}

func (d *DistributorSubscriptionService) InitPayment(ctx context.Context, distId int) (SubscribeResponse, error) {
	//get current subscription
	//init subscription checkout and return
	return SubscribeResponse{}, nil
}

func (d *DistributorSubscriptionService) GetSubscriptions(ctx context.Context) (GetAllSubscriptionPlanResponse, error) {
	return GetAllSubscriptionPlanResponse{}, nil
}

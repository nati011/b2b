package distributor_subscription

import (
	"context"

	"b2b.nati011.github.com/internal/core/domain/distributor"
	port "b2b.nati011.github.com/internal/port/domain/distributor_subscription"
)

type SubscribeResponse struct {
	Id          int    `json:"id"`
	CheckoutUrl string `json:"checkout_url"`
	TxRef       string `json:"tx_ref"`
}

type GetSubscriptionPlanResponse struct {
}

type CreatePlanRequest struct {
}

// TODO: acitvate and deactivate plans
type Prodvider interface {
	GetPlan(ctx context.Context) (GetSubscriptionPlanResponse, error)
	CreatePlan(ctx context.Context, req *CreatePlanRequest) (int, error)
	Place(ctx context.Context) (SubscribeResponse, error)
	InitPayment(ctx context.Context, distId int) (SubscribeResponse, error)
}

type DistributorSubscriptionService struct {
	DB port.DB
	DS distributor.Provider
}

func (d *DistributorSubscriptionService) GetPlans(ctx context.Context) (GetSubscriptionPlanResponse, error) {
	return GetSubscriptionPlanResponse{}, nil
}

func (d *DistributorSubscriptionService) Place(ctx context.Context) (SubscribeResponse, error) {
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

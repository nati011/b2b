package distributorApproval

import (
	"context"
)

type PlaceRequest struct {
	SubscriptionPlanId int
	DistributorId      int
	PaymentPartnerId   int
	Status             string
}

type GetSubscriptionResponse struct {
	Id                 int
	SubscriptionPlanId int
	DistributorId      int
	PaymentPartnerId   int
	Status             string
}

type UpdateSubscriptionRequest struct {
	Id     int
	Status string
}

type GetAllSubscriptionResponse struct {
	List []GetSubscriptionResponse
}

type CreatePlanRequest struct {
	Name        string
	Price       float64
	TermInMonth int
	Description string
}

type GetPlanResponse struct {
	Id          int
	Name        string
	Price       float64
	TermInMonth int
	Description string
}

type GetAllPlanResponse struct {
	List []GetPlanResponse
}

type Reader interface {
	GetPlan(ctx context.Context, id int) (GetPlanResponse, error)
	GetAllPlan(ctx context.Context) (GetAllPlanResponse, error)
	GetAllSubscriptions(ctx context.Context) (GetAllSubscriptionResponse, error)
	GetSubscription(ctx context.Context, subId int) (GetSubscriptionResponse, error)
	GetSubscriptionByDistributorId(ctx context.Context, distId int) (GetSubscriptionResponse, error)
}

type Writer interface {
	Place(ctx context.Context, req *PlaceRequest) (int, error)
	CreatePlan(ctx context.Context, req *CreatePlanRequest) (int, error)
	UpdateStatus(ctx context.Context, req *UpdateSubscriptionRequest) error
}

type DB interface {
	Reader
	Writer
}

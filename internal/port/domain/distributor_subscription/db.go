package distributorApproval

import (
	"context"
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
	GetPlans(ctx context.Context) (GetAllPlanResponse, error)
	GetAllSubscriptionResponse(ctx context.Context) (GetAllSubscriptionResponse, error)
}

type Writer interface {
	Place(ctx context.Context, req *PlaceRequest) error
	CreatePlan(ctx context.Context, req *CreatePlanRequest) (int, error)
}

type DB interface {
	Reader
	Writer
}

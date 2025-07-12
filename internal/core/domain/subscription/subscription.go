package subscription

import "context"

type CreatePlanRequest struct {
}

type GetPlanResponse struct {
}

type Provider interface {
	CreatePlan(ctx context.Context, req *CreatePlanRequest) (int, error)
}

type SubscriptionService struct {
}

func NewSubscriptionService() Provider {
	return &SubscriptionService{}
}

func (s *SubscriptionService) CreatePlan(ctx context.Context, req *CreatePlanRequest) (int, error) {
	return 0, nil
}

func (s *SubscriptionService) GetPlan(ctx context.Context, id int) (GetPlanResponse, error) {
	return GetPlanResponse{}, nil
}

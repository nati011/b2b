package distributorsubscription

import (
	"context"

	port "b2b.nati011.github.com/internal/port/domain/distributor_subscription"
)

type MockPlan struct {
	Id          int
	Name        string
	Price       float64
	TermInMonth int
	Description string
}

type MockSubscription struct {
	Id                 int
	SubscriptionPlanId int
	DistributorId      int
	Status             string
}

type Mock struct {
	Plans []MockPlan
	Subs  []MockSubscription
}

func NewMock() port.DB {
	return &Mock{}
}

func (m *Mock) CreatePlan(ctx context.Context, req *port.CreatePlanRequest) (int, error) {
	newId := len(m.Plans) + 1
	m.Plans = append(m.Plans, MockPlan{
		Id:          newId,
		Name:        req.Name,
		Price:       req.Price,
		TermInMonth: req.TermInMonth,
		Description: req.Description,
	})
	return newId, nil
}

func (m *Mock) GetPlans(ctx context.Context) (port.GetAllPlanResponse, error) {
	response := port.GetAllPlanResponse{}
	for _, i := range m.Plans {
		response.List = append(response.List, port.GetPlanResponse(i))
	}
	return port.GetAllPlanResponse{}, nil
}

func (m *Mock) Place(ctx context.Context, req *port.PlaceRequest) error {
	newId := len(m.Subs) + 1
	m.Subs = append(m.Subs, MockSubscription{
		Id:                 newId,
		SubscriptionPlanId: req.SubscriptionPlanId,
		DistributorId:      req.DistributorId,
		Status:             req.Status,
	})
	return nil
}

func (m *Mock) GetAllSubscriptions(ctx context.Context) (port.GetAllSubscriptionResponse, error) {
	response := port.GetAllSubscriptionResponse{}
	for _, i := range m.Subs {
		response.List = append(response.List, port.GetSubscriptionResponse(i))
	}
	return port.GetAllSubscriptionResponse{}, nil
}

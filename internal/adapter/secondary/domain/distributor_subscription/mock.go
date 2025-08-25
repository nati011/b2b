package distributor_subscription

import (
	"context"

	port_commons "b2b.nati011.github.com/internal/port/commons/db"
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

type DistributorSubscriptionMock struct {
	Plans []MockPlan
	Subs  []MockSubscription
}

func NewDistributorSubscriptionMock() port.DB {
	return &DistributorSubscriptionMock{}
}

func (m *DistributorSubscriptionMock) CreatePlan(ctx context.Context, req *port.CreatePlanRequest) (int, error) {
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

func (m *DistributorSubscriptionMock) GetPlans(ctx context.Context) (port.GetAllPlanResponse, error) {
	response := port.GetAllPlanResponse{}
	for _, i := range m.Plans {
		response.List = append(response.List, port.GetPlanResponse(i))
	}
	if len(response.List) == 0 {
		return port.GetAllPlanResponse{}, port_commons.ErrSysNoRows
	}
	return response, nil
}

func (m *DistributorSubscriptionMock) Place(ctx context.Context, req *port.PlaceRequest) (int, error) {
	newId := len(m.Subs) + 1
	m.Subs = append(m.Subs, MockSubscription{
		Id:                 newId,
		SubscriptionPlanId: req.SubscriptionPlanId,
		DistributorId:      req.DistributorId,
		Status:             req.Status,
	})
	return newId, nil
}

func (m *DistributorSubscriptionMock) GetAllSubscriptions(ctx context.Context) (port.GetAllSubscriptionResponse, error) {
	response := port.GetAllSubscriptionResponse{}
	for _, i := range m.Subs {
		response.List = append(response.List, port.GetSubscriptionResponse(i))
	}
	return port.GetAllSubscriptionResponse{}, nil
}

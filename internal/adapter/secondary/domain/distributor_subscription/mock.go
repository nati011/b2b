package db

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
	PaymentPartnerId   int
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

func (m *DistributorSubscriptionMock) GetAllPlan(ctx context.Context) (port.GetAllPlanResponse, error) {
	response := port.GetAllPlanResponse{}
	for _, i := range m.Plans {
		response.List = append(response.List, port.GetPlanResponse(i))
	}
	if len(response.List) == 0 {
		return port.GetAllPlanResponse{}, port_commons.ErrSysNoRows
	}
	return response, nil
}

func (m *DistributorSubscriptionMock) GetPlan(ctx context.Context, id int) (port.GetPlanResponse, error) {
	for _, i := range m.Plans {
		if i.Id == id {
			return port.GetPlanResponse(i), nil
		}
	}
	return port.GetPlanResponse{}, port_commons.ErrSysNoRows
}

func (m *DistributorSubscriptionMock) Place(ctx context.Context, req *port.PlaceRequest) (int, error) {
	newId := len(m.Subs) + 1
	m.Subs = append(m.Subs, MockSubscription{
		Id:                 newId,
		SubscriptionPlanId: req.SubscriptionPlanId,
		DistributorId:      req.DistributorId,
		PaymentPartnerId:   req.PaymentPartnerId,
		Status:             req.Status,
	})
	return newId, nil
}

func (m *DistributorSubscriptionMock) GetAllSubscriptions(ctx context.Context) (port.GetAllSubscriptionResponse, error) {
	response := port.GetAllSubscriptionResponse{}
	for _, i := range m.Subs {
		response.List = append(response.List, port.GetSubscriptionResponse(i))
	}
	if len(response.List) == 0 {
		return port.GetAllSubscriptionResponse{}, port_commons.ErrSysNoRows
	}
	return response, nil
}

func (m *DistributorSubscriptionMock) GetSubscription(ctx context.Context, subId int) (port.GetSubscriptionResponse, error) {
	for _, i := range m.Subs {
		if i.Id == subId {
			return port.GetSubscriptionResponse(i), nil
		}
	}
	return port.GetSubscriptionResponse{}, port_commons.ErrSysNoRows
}

func (m *DistributorSubscriptionMock) UpdateStatus(ctx context.Context, req *port.UpdateSubscriptionRequest) error {
	subs := []MockSubscription{}
	for _, i := range m.Subs {
		if i.Id == req.Id {
			subs = append(subs, MockSubscription{
				Id:                 i.Id,
				SubscriptionPlanId: i.SubscriptionPlanId,
				DistributorId:      i.DistributorId,
				PaymentPartnerId:   i.PaymentPartnerId,
				Status:             req.Status,
			})
		} else {
			subs = append(subs, MockSubscription(i))
		}

	}
	m.Subs = subs
	return nil
}

func (m *DistributorSubscriptionMock) GetSubscriptionByDistributorId(ctx context.Context, distributorId int) (port.GetSubscriptionResponse, error) {
	for _, i := range m.Subs {
		if i.DistributorId == distributorId {
			return port.GetSubscriptionResponse(i), nil
		}

	}
	return port.GetSubscriptionResponse{}, port_commons.ErrSysNoRows
}

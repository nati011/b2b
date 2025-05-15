package db

import (
	"context"

	port_commons "b2b.nati011.github.com/internal/port/commons/db"
	port "b2b.nati011.github.com/internal/port/domain/distributor"
)

type MockDistributor struct {
	Id          int
	Name        string
	Tin         string
	Latitude    string
	Longitude   string
	GeneralZone string
	Region      string
	Woreda      string
	IsActive    bool
}

type MockUserAgent struct {
	Id int
}

type Mock struct {
	distributors []MockDistributor
	userAgents   []MockUserAgent
}

func NewMock() port.DB {
	return &Mock{}
}

func (m *Mock) Create(ctx context.Context, req port.CreateRequest) (int, error) {
	newId := len(m.distributors) + 1
	m.distributors = append(m.distributors, MockDistributor{
		Id:          newId,
		Name:        req.Name,
		Tin:         req.Tin,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		GeneralZone: req.GeneralZone,
		Region:      req.Region,
		Woreda:      req.Woreda,
	})

	user_id, err := m.CreateDistributorUser(ctx, &port.CreateUserAgentRequest{
		User_id: req.UserId,
	})
	if err != nil {
		return 0, port_commons.ErrSysUnknown
	}
	m.userAgents = append(m.userAgents, MockUserAgent{
		Id: user_id,
	})
	return newId, nil
}

func (m Mock) CreateDistributorUser(ctx context.Context, req *port.CreateUserAgentRequest) (int, error) {
	m.userAgents = append(m.userAgents, MockUserAgent{
		Id: req.User_id,
	})
	return req.User_id, nil
}

func (m *Mock) UpdateName(ctx context.Context, req *port.UpdateNameRequest) error {
	distributors := []MockDistributor{}
	for _, i := range m.distributors {
		if i.Id == req.Id {
			distributors = append(distributors, MockDistributor{
				Id:          i.Id,
				Name:        req.Name,
				Tin:         i.Tin,
				Latitude:    i.Latitude,
				Longitude:   i.Longitude,
				GeneralZone: i.GeneralZone,
				Region:      i.Region,
				Woreda:      i.Woreda,
				IsActive:    i.IsActive,
			})
		} else {
			distributors = append(distributors, MockDistributor{
				Id:          i.Id,
				Name:        i.Name,
				Tin:         i.Tin,
				Latitude:    i.Latitude,
				Longitude:   i.Longitude,
				GeneralZone: i.GeneralZone,
				Region:      i.Region,
				Woreda:      i.Woreda,
				IsActive:    i.IsActive,
			})
		}
	}
	m.distributors = distributors
	return nil
}

func (m *Mock) UpdateTin(ctx context.Context, req *port.UpdateTinRequest) error {
	distributors := []MockDistributor{}
	for _, i := range m.distributors {
		if i.Id == req.Id {
			distributors = append(distributors, MockDistributor{
				Id:          i.Id,
				Name:        i.Name,
				Tin:         req.Tin,
				Latitude:    i.Latitude,
				Longitude:   i.Longitude,
				GeneralZone: i.GeneralZone,
				Region:      i.Region,
				Woreda:      i.Woreda,
				IsActive:    i.IsActive,
			})
		} else {
			distributors = append(distributors, MockDistributor{
				Id:          i.Id,
				Name:        i.Name,
				Tin:         i.Tin,
				Latitude:    i.Latitude,
				Longitude:   i.Longitude,
				GeneralZone: i.GeneralZone,
				Region:      i.Region,
				Woreda:      i.Woreda,
				IsActive:    i.IsActive,
			})
		}
	}
	m.distributors = distributors
	return nil
}

func (m *Mock) Get(ctx context.Context, id int) (port.GetResponse, error) {
	for _, i := range m.distributors {
		if i.Id == id {
			return port.GetResponse{
				Id:          i.Id,
				Name:        i.Name,
				Tin:         i.Tin,
				Latitude:    i.Latitude,
				Longitude:   i.Longitude,
				GeneralZone: i.GeneralZone,
				Region:      i.Region,
				Woreda:      i.Woreda,
				IsActive:    i.IsActive,
			}, nil
		}
	}
	return port.GetResponse{}, port_commons.ErrSysNoRows
}

func (m *Mock) GetAll(ctx context.Context) (port.GetAllResponse, error) {
	resp := []port.GetResponse{}
	for _, i := range m.distributors {
		resp = append(resp, port.GetResponse{
			Id:          i.Id,
			Name:        i.Name,
			Tin:         i.Tin,
			Latitude:    i.Latitude,
			Longitude:   i.Longitude,
			GeneralZone: i.GeneralZone,
			Region:      i.Region,
			Woreda:      i.Woreda,
			IsActive:    i.IsActive,
		})
	}
	if len(resp) == 0 {
		return port.GetAllResponse{}, port_commons.ErrSysNoRows
	}
	return port.GetAllResponse{
		List: resp,
	}, nil
}

func (m *Mock) GetByName(ctx context.Context, name string) (port.GetAllResponse, error) {
	resp := []port.GetResponse{}
	for _, i := range m.distributors {
		if i.Name == name {
			resp = append(resp, port.GetResponse{
				Id:          i.Id,
				Name:        i.Name,
				Tin:         i.Tin,
				Latitude:    i.Latitude,
				Longitude:   i.Longitude,
				GeneralZone: i.GeneralZone,
				Region:      i.Region,
				Woreda:      i.Woreda,
				IsActive:    i.IsActive,
			})
		}
	}

	if len(resp) == 0 {
		return port.GetAllResponse{}, port_commons.ErrSysNoRows
	}
	return port.GetAllResponse{
		List: resp,
	}, nil
}

func (m *Mock) GetByTin(ctx context.Context, tin string) (port.GetResponse, error) {
	for _, i := range m.distributors {
		if i.Tin == tin {
			return port.GetResponse{
				Id:          i.Id,
				Name:        i.Name,
				Tin:         i.Tin,
				Latitude:    i.Latitude,
				Longitude:   i.Longitude,
				GeneralZone: i.GeneralZone,
				Region:      i.Region,
				Woreda:      i.Woreda,
				IsActive:    i.IsActive,
			}, nil
		}
	}
	return port.GetResponse{}, port_commons.ErrSysNoRows
}

func (m *Mock) GetAllUserAgents(ctx context.Context, id int) (port.GetAllUserResponse, error) {
	var resp = port.GetAllUserResponse{}
	for _, i := range m.userAgents {
		resp.List = append(resp.List, port.GetUserResponse{
			Id: i.Id,
		})
	}
	if len(resp.List) == 0 {
		return port.GetAllUserResponse{}, port_commons.ErrSysNoRows
	}
	return resp, nil
}

func (m *Mock) Activate(ctx context.Context, id int) error {
	distributors := []MockDistributor{}
	for _, i := range m.distributors {
		if i.Id != id {
			distributors = append(distributors, MockDistributor{
				Id:          i.Id,
				Name:        i.Name,
				Tin:         i.Tin,
				Latitude:    i.Latitude,
				Longitude:   i.Longitude,
				GeneralZone: i.GeneralZone,
				Region:      i.Region,
				Woreda:      i.Woreda,
				IsActive:    i.IsActive,
			})
		} else {
			distributors = append(distributors, MockDistributor{
				Id:          i.Id,
				Name:        i.Name,
				Tin:         i.Tin,
				Latitude:    i.Latitude,
				Longitude:   i.Longitude,
				GeneralZone: i.GeneralZone,
				Region:      i.Region,
				Woreda:      i.Woreda,
				IsActive:    true,
			})
		}
	}
	m.distributors = distributors
	return nil
}
func (m *Mock) Dectivate(ctx context.Context, id int) error {
	distributors := []MockDistributor{}
	for _, i := range m.distributors {
		if i.Id != id {
			distributors = append(distributors, MockDistributor{
				Id:          i.Id,
				Name:        i.Name,
				Tin:         i.Tin,
				Latitude:    i.Latitude,
				Longitude:   i.Longitude,
				GeneralZone: i.GeneralZone,
				Region:      i.Region,
				Woreda:      i.Woreda,
				IsActive:    i.IsActive,
			})
		} else {
			distributors = append(distributors, MockDistributor{
				Id:          i.Id,
				Name:        i.Name,
				Tin:         i.Tin,
				Latitude:    i.Latitude,
				Longitude:   i.Longitude,
				GeneralZone: i.GeneralZone,
				Region:      i.Region,
				Woreda:      i.Woreda,
				IsActive:    false,
			})
		}
	}
	m.distributors = distributors
	return nil
}

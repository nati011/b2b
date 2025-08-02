package db

import (
	"context"
	"log"

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

type MockDistributorUser struct {
	UserId        int
	DistributorId int
}

type MockUserAgent struct {
	Id        int
	FirstName string
	LastName  string
	Email     string
	Phone     string
	Username  string
}

type Mock struct {
	distributors    []MockDistributor
	userAgents      []MockUserAgent
	distributorUser []MockDistributorUser
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

	m.userAgents = append(m.userAgents, MockUserAgent{
		Id: req.UserId,
	})
	m.distributorUser = append(m.distributorUser, MockDistributorUser{
		UserId:        req.UserId,
		DistributorId: newId,
	})
	return newId, nil
}

func (m *Mock) Remove(ctx context.Context, id int) error {
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
		}
	}
	m.distributors = distributors
	return nil
}

func (m *Mock) CreateDistributorUser(ctx context.Context, req *port.CreateUserAgentRequest) (int, error) {
	m.userAgents = append(m.userAgents, MockUserAgent{
		Id: req.User_id,
	})

	log.Printf("User Agents: %v", m.userAgents)
	if m.distributorUser == nil {
		m.distributorUser = make([]MockDistributorUser, 0)
	}
	m.distributorUser = append(m.distributorUser, MockDistributorUser{
		UserId:        req.User_id,
		DistributorId: req.Distributor_Id,
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

func (m *Mock) GetAllUserDetail(ctx context.Context, distributor_id int) (port.GetAllUserDetailResponse, error) {
	var distributorUsers []int
	var resp []port.GetUserDetailResponse

	for _, i := range m.distributorUser {
		if i.DistributorId == distributor_id {
			resp = append(
				resp, port.GetUserDetailResponse{
					Id: i.UserId,
				},
			)
		}
	}

	if len(distributorUsers) == 0 {
		return port.GetAllUserDetailResponse{}, port_commons.ErrSysNoRows
	}

	return port.GetAllUserDetailResponse{
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

func (m *Mock) GetByStatus(ctx context.Context, status bool) (port.GetAllResponse, error) {
	resp := []port.GetResponse{}
	for _, i := range m.distributors {
		if i.IsActive == status {
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

func (m *Mock) GetByApprovalStatus(ctx context.Context, status string) (port.GetAllResponse, error) {
	return port.GetAllResponse{}, nil
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
func (m *Mock) GetByUserId(ctx context.Context, user_id int) (port.GetResponse, error) {
	var DistributorId int
	for _, i := range m.distributorUser {
		if i.UserId == user_id {
			DistributorId = i.DistributorId
		}
	}
	for _, i := range m.distributors {
		if i.Id == DistributorId {
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

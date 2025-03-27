package adapter

import (
	"context"

	port "b2b.nati011.github.com/internal/port/application/retailer"
)

type MockRetailer struct {
	Id          int
	Name        string
	Tin         string
	Latitude    string
	Longitude   string
	GeneralZone string
	Region      string
	Woreda      string
}

type MockUserAgent struct {
	Id int
}

type Mock struct {
	retailers  []MockRetailer
	userAgents []MockUserAgent
}

func NewMock() port.DB {
	return &Mock{}
}

func (m *Mock) Create(ctx context.Context, req port.CreateRequest) (int, error) {
	newId := len(m.retailers) + 1
	m.retailers = append(m.retailers, MockRetailer{
		Id:          newId,
		Name:        req.Name,
		Tin:         req.Tin,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		GeneralZone: req.GeneralZone,
		Region:      req.Region,
		Woreda:      req.Woreda,
	})

	err := m.CreateRetailerUser(ctx, &port.CreateUserAgentRequest{
		User_id: req.UserId,
	})
	if err != nil {
		return 0, port.ErrSysUnknown
	}

	return newId, nil
}

func (m Mock) CreateRetailerUser(ctx context.Context, req *port.CreateUserAgentRequest) error {
	m.userAgents = append(m.userAgents, MockUserAgent{
		Id: req.User_id,
	})
	return nil
}

func (m *Mock) UpdateName(ctx context.Context, req *port.UpdateNameRequest) error {
	retailers := []MockRetailer{}
	for _, i := range m.retailers {
		if i.Id == req.Id {
			retailers = append(retailers, MockRetailer{
				Id:          i.Id,
				Name:        req.Name,
				Tin:         i.Tin,
				Latitude:    i.Latitude,
				Longitude:   i.Longitude,
				GeneralZone: i.GeneralZone,
				Region:      i.Region,
				Woreda:      i.Woreda,
			})
		} else {
			retailers = append(retailers, MockRetailer{
				Id:          i.Id,
				Name:        i.Name,
				Tin:         i.Tin,
				Latitude:    i.Latitude,
				Longitude:   i.Longitude,
				GeneralZone: i.GeneralZone,
				Region:      i.Region,
				Woreda:      i.Woreda,
			})
		}
	}
	m.retailers = retailers
	return nil
}

func (m *Mock) UpdateTin(ctx context.Context, req *port.UpdateTinRequest) error {
	retailers := []MockRetailer{}
	for _, i := range m.retailers {
		if i.Id == req.Id {
			retailers = append(retailers, MockRetailer{
				Id:          i.Id,
				Name:        i.Name,
				Tin:         req.Tin,
				Latitude:    i.Latitude,
				Longitude:   i.Longitude,
				GeneralZone: i.GeneralZone,
				Region:      i.Region,
				Woreda:      i.Woreda,
			})
		} else {
			retailers = append(retailers, MockRetailer{
				Id:          i.Id,
				Name:        i.Name,
				Tin:         i.Tin,
				Latitude:    i.Latitude,
				Longitude:   i.Longitude,
				GeneralZone: i.GeneralZone,
				Region:      i.Region,
				Woreda:      i.Woreda,
			})
		}
	}
	m.retailers = retailers
	return nil
}

func (m *Mock) Get(ctx context.Context, id int) (port.GetResponse, error) {
	for _, i := range m.retailers {
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
			}, nil
		}
	}
	return port.GetResponse{}, port.ErrSysNoRows
}

func (m *Mock) GetAll(ctx context.Context) (port.GetAllResponse, error) {
	resp := []port.GetResponse{}
	for _, i := range m.retailers {
		resp = append(resp, port.GetResponse{
			Id:          i.Id,
			Name:        i.Name,
			Tin:         i.Tin,
			Latitude:    i.Latitude,
			Longitude:   i.Longitude,
			GeneralZone: i.GeneralZone,
			Region:      i.Region,
			Woreda:      i.Woreda,
		})
	}
	if len(resp) == 0 {
		return port.GetAllResponse{}, port.ErrSysNoRows
	}
	return port.GetAllResponse{
		List: resp,
	}, nil
}

func (m *Mock) GetByName(ctx context.Context, name string) (port.GetAllResponse, error) {
	resp := []port.GetResponse{}
	for _, i := range m.retailers {
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
			})
		}
	}

	if len(resp) == 0 {
		return port.GetAllResponse{}, port.ErrSysNoRows
	}
	return port.GetAllResponse{
		List: resp,
	}, nil
}

func (m *Mock) GetByTin(ctx context.Context, tin string) (port.GetResponse, error) {
	for _, i := range m.retailers {
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
			}, nil
		}
	}
	return port.GetResponse{}, port.ErrSysNoRows
}

func (m *Mock) GetAllUserAgents(ctx context.Context, id int) (port.GetAllUserResponse, error) {
	var resp = port.GetAllUserResponse{}
	for _, i := range m.userAgents {
		resp.List = append(resp.List, port.GetUserResponse{
			Id: i.Id,
		})
	}
	if len(resp.List) == 0 {
		return port.GetAllUserResponse{}, port.ErrSysNoRows
	}
	return resp, nil
}

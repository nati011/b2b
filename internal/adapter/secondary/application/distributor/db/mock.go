package db

import (
	"context"

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
}

type MockUserAgent struct {
	Id int
}

type Mock struct {
<<<<<<< HEAD
	retailers  []MockDistributor
	userAgents []MockUserAgent
=======
	Distributors []MockDistributor
	Businesses   []MockBusiness
}

func (m *Mock) GetById(ctx context.Context, id int) (port.GetResponse, error) {
	resp := port.GetResponse{}
	for _, i := range m.Distributors {
		resp = port.GetResponse{
			Id: i.Id,
		}
	}
	if (resp == port.GetResponse{}) {
		return port.GetResponse{}, port.ErrSysNoRows
	}

	return resp, nil
}

func (m *Mock) Create(ctx context.Context, req *port.RegisterDistributorRequest) (int, error) {
	distributorId := rand.Int()
	m.Distributors = append(m.Distributors, MockDistributor{
		Id: distributorId,
	})
	return distributorId, nil

}

func (m *Mock) GetAll(ctx context.Context) (port.GetAllResponse, error) {
	resp := []port.GetResponse{}
	for _, i := range m.Distributors {

		resp = append(resp, port.GetResponse{
			Id: i.Id,
		})
	}
	if len(resp) == 0 {
		return port.GetAllResponse{}, port.ErrSysNoRows
	}

	return port.GetAllResponse{
		List: resp,
	}, nil
}

func (m *Mock) GetBusinessById(ctx context.Context, id int) (resp port.GetBusinessResponse, err error) {
	for _, i := range m.Businesses {
		log.Print(i)
		resp = port.GetBusinessResponse{
			Id: i.Id,
		}
	}

	if (resp == port.GetBusinessResponse{}) {
		return port.GetBusinessResponse{}, port.ErrSysNoRows
	}

	return resp, nil
}

func (m *Mock) GetByDistributorId(ctx context.Context, distributorId int) (port.GetBusinessResponse, error) {
	resp := port.GetBusinessResponse{}
	for _, i := range m.Businesses {
		resp = port.GetBusinessResponse{
			Id: i.DistributorId,
		}
	}
	if (resp == port.GetBusinessResponse{}) {
		return port.GetBusinessResponse{}, port.ErrSysNoRows
	}

	return resp, nil
}

func (m *Mock) GetBusinessAll(ctx context.Context) (port.GetAllResponse, error) {
	resp := []port.GetResponse{}
	for _, i := range m.Businesses {

		resp = append(resp, port.GetResponse{
			Id: i.Id,
		})
	}
	if len(resp) == 0 {
		return port.GetAllResponse{}, port.ErrSysNoRows
	}

	return port.GetAllResponse{
		List: resp,
	}, nil
}

func (m *Mock) CreateBusiness(ctx context.Context, req *port.CreateBusinessInformation) (port.CreateBusinessResponse, error) {
	businessId := rand.Int()
	m.Businesses = append(m.Businesses, MockBusiness{
		Id:            businessId,
		Name:          req.Name,
		Tin:           req.Tin,
		DistributorId: req.DistributorId,
	})

	response := port.CreateBusinessResponse{
		BusinessId: businessId,
	}
	return response, nil
}

func (m *Mock) UpdateBusiness(ctx context.Context, req *port.UpdateBusinessRequest) (int, error) {
	updatedDistributor := []MockBusiness{}
	var updatedResourceId int
	for _, i := range m.Businesses {
		if req.Id == i.Id {
			updatedResourceId = i.Id
			updatedDistributor = append(updatedDistributor, MockBusiness{
				Id:   req.Id,
				Name: req.Name,
				Tin:  req.Tin,
			})
		} else {
			updatedDistributor = append(updatedDistributor, MockBusiness{
				Id:   i.Id,
				Name: i.Name,
				Tin:  i.Tin,
			})
		}

	}
	m.Businesses = updatedDistributor
	return updatedResourceId, nil

>>>>>>> 8f0b9404 (init distributor refactor)
}

func NewMock() port.DB {
	return &Mock{}
}

func (m *Mock) Create(ctx context.Context, req port.CreateRequest) (int, error) {
	newId := len(m.retailers) + 1
	m.retailers = append(m.retailers, MockDistributor{
		Id:          newId,
		Name:        req.Name,
		Tin:         req.Tin,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		GeneralZone: req.GeneralZone,
		Region:      req.Region,
		Woreda:      req.Woreda,
	})

	_, err := m.CreateDistributorUser(ctx, &port.CreateUserAgentRequest{
		User_id: req.UserId,
	})
	if err != nil {
		return 0, port.ErrSysUnknown
	}

	return newId, nil
}

func (m Mock) CreateDistributorUser(ctx context.Context, req *port.CreateUserAgentRequest) (int, error) {
	m.userAgents = append(m.userAgents, MockUserAgent{
		Id: req.User_id,
	})
	return req.User_id, nil
}

func (m *Mock) UpdateName(ctx context.Context, req *port.UpdateNameRequest) error {
	retailers := []MockDistributor{}
	for _, i := range m.retailers {
		if i.Id == req.Id {
			retailers = append(retailers, MockDistributor{
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
			retailers = append(retailers, MockDistributor{
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
	retailers := []MockDistributor{}
	for _, i := range m.retailers {
		if i.Id == req.Id {
			retailers = append(retailers, MockDistributor{
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
			retailers = append(retailers, MockDistributor{
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

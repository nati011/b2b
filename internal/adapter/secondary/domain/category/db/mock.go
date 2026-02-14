package category

import (
	"context"
	"math/rand"
	"time"

	port_commons "b2b.nati011.github.com/internal/port/commons/db"
	port "b2b.nati011.github.com/internal/port/domain/category"
)

type MockCategory struct {
	Id   int
	Name string
}

type Mock struct {
	Categories []MockCategory
}

func NewMock() port.DB {
	return &Mock{}
}

func (m *Mock) Create(ctx context.Context, req *port.CreateRequest) (int, error) {
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator
	id := rand.Intn(1000-10+1) + 10
	m.Categories = append(m.Categories, MockCategory{
		Id:   id,
		Name: req.Name,
	})
	return id, nil
}

func (m *Mock) GetAll(ctx context.Context) (port.GetAllResponse, error) {
	resp := port.GetAllResponse{}
	for _, i := range m.Categories {
		resp.List = append(resp.List, port.GetResponse(i))
	}
	if len(resp.List) == 0 {
		return resp, port_commons.ErrSysNoRows
	}
	return resp, nil
}

func (m *Mock) Get(ctx context.Context, id int) (port.GetResponse, error) {
	for _, i := range m.Categories {
		if i.Id == id {
			return port.GetResponse(i), nil
		}
	}
	return port.GetResponse{}, port_commons.ErrSysNoRows
}

func (m *Mock) Remove(ctx context.Context, id int) error {
	resp := []MockCategory{}
	for _, i := range m.Categories {
		if i.Id != id {
			resp = append(resp, i)
		}
	}
	m.Categories = resp
	return nil
}

func (m *Mock) Update(ctx context.Context, req *port.UpdateRequest) error {
	resp := []MockCategory{}
	for _, i := range m.Categories {
		if i.Id == req.Id {
			resp = append(resp, MockCategory{
				Id:   i.Id,
				Name: req.Name,
			})
		} else {
			resp = append(resp, i)
		}
	}
	m.Categories = resp
	return nil
}

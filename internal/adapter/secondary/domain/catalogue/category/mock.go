package category

import (
	"context"
	"math/rand"
	"time"

	port "b2b.nati011.github.com/internal/port/catalogue/category"
)

type MockCategory struct {
	Id   int
	Name string
	Desc string
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
		Desc: req.Desc,
	})
	return id, nil
}

func (m *Mock) GetAll(ctx context.Context) (port.GetAllResponse, error) {
	resp := port.GetAllResponse{}
	for _, i := range m.Categories {
		resp.List = append(resp.List, port.GetResponse(i))
	}
	if len(resp.List) == 0 {
		return resp, port.ErrSysNoRows
	}
	return resp, nil
}

func (m *Mock) Get(ctx context.Context, id int) (port.GetResponse, error) {
	resp := port.GetResponse{}
	for _, i := range m.Categories {
		resp = port.GetResponse(i)
	}
	if resp.Id != id {
		return resp, port.ErrSysNoRows
	}
	return resp, nil
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

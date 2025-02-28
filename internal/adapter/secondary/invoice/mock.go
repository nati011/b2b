package invoice

import (
	"context"
	"math/rand"
	"time"

	port "b2b.nati011.github.com/internal/port/invoice"
)

type MockInvoice struct {
	Id           int
	Created_Date time.Time
	ExternalId   string
	Status       string
}

type Mock struct {
	invoices []MockInvoice
}

func NewMock() port.DB {
	return &Mock{}
}

func (m *Mock) Get(ctx context.Context, id int) (port.GetResponse, error) {
	for _, i := range m.invoices {
		if i.Id == id {
			return port.GetResponse(i), nil
		}
	}
	return port.GetResponse{}, port.ErrSysNoRows
}

func (m *Mock) GetAll(ctx context.Context) (port.GetAllResponse, error) {
	res := port.GetAllResponse{}
	for _, i := range m.invoices {
		res.List = append(res.List, port.GetResponse(i))

	}
	if len(res.List) == 0 {
		return port.GetAllResponse{}, port.ErrSysNoRows
	}
	return res, nil
}

func (m *Mock) GetByExternalId(ctx context.Context, extId string) (port.GetAllResponse, error) {
	res := port.GetAllResponse{}
	for _, i := range m.invoices {
		if i.ExternalId == extId {
			res.List = append(res.List, port.GetResponse(i))
		}
	}
	if len(res.List) == 0 {
		return port.GetAllResponse{}, port.ErrSysNoRows
	}
	return res, nil
}

func (m *Mock) GetByStatus(ctx context.Context, status string) (port.GetAllResponse, error) {
	res := port.GetAllResponse{}
	for _, i := range m.invoices {
		if i.Status == status {
			res.List = append(res.List, port.GetResponse(i))
		}
	}
	if len(res.List) == 0 {
		return port.GetAllResponse{}, port.ErrSysNoRows
	}
	return res, nil
}

func (m *Mock) Create(ctx context.Context, req *port.CreateRequest) (int, error) {
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator
	id := rand.Intn(1000-10+1) + 10
	m.invoices = append(m.invoices, MockInvoice{
		Id:           id,
		Created_Date: time.Now(),
		ExternalId:   req.ExternalId,
		Status:       req.Status,
	})
	return id, nil
}

func (m *Mock) UpdateExternalId(ctx context.Context, req *port.UpdateExternalIdRequest) error {
	res := []MockInvoice{}
	for _, i := range m.invoices {
		if i.Id == req.Id {
			res = append(res, MockInvoice{
				Id:           i.Id,
				Created_Date: i.Created_Date,
				ExternalId:   req.ExternalId,
				Status:       i.Status,
			})
		} else {
			res = append(res, MockInvoice(i))
		}
	}
	m.invoices = res
	return nil
}

func (m *Mock) UpdateStatus(ctx context.Context, req *port.UpdateStatusRequest) error {
	res := []MockInvoice{}
	for _, i := range m.invoices {
		if i.Id == req.Id {
			res = append(res, MockInvoice{
				Id:           i.Id,
				Created_Date: i.Created_Date,
				ExternalId:   i.ExternalId,
				Status:       req.Status,
			})
		} else {
			res = append(res, MockInvoice(i))
		}
	}
	m.invoices = res
	return nil
}

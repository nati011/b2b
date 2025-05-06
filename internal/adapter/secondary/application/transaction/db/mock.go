package transaction

import (
	"context"
	"time"

	port "b2b.nati011.github.com/internal/port/application/transaction/db"
	port_commons "b2b.nati011.github.com/internal/port/commons/db"
)

type GetResponse struct {
	Id        int
	Date      time.Time
	Amount    float64
	PartnerId int
	Status    string
}

type GetAllResponse struct {
	List []GetAllResponse
}

type CreateRequest struct {
	Amount    int64
	PartnerId int
	Status    string
	TxRef     string
}

type GetByParamRequest struct {
	Date      time.Time
	PartnerId int
	TxRef     string
	Status    string
}

type MockTransaction struct {
	Id        int
	Date      time.Time
	Amount    float64
	PartnerId int
	TxRef     string
	Status    string
}

type Mock struct {
	resources []MockTransaction
}

func NewMock() port.DB {
	return &Mock{}
}

func (m *Mock) GetByID(ctx context.Context, id int) (port.GetResponse, error) {
	for _, i := range m.resources {
		if i.Id == id {
			return port.GetResponse{
				Id:        i.Id,
				Date:      i.Date,
				Amount:    i.Amount,
				PartnerId: i.PartnerId,
				TxRef:     i.TxRef,
				Status:    i.Status,
			}, nil
		}
	}
	return port.GetResponse{}, port_commons.ErrSysNoRows
}

func (m *Mock) GetAll(ctx context.Context) (port.GetAllResponse, error) {
	response := []port.GetResponse{}
	for _, i := range m.resources {
		response = append(response, port.GetResponse{
			Id:        i.Id,
			Date:      i.Date,
			Amount:    i.Amount,
			PartnerId: i.PartnerId,
			TxRef:     i.TxRef,
			Status:    i.Status,
		})
	}
	if len(response) == 0 {
		return port.GetAllResponse{}, port_commons.ErrSysNoRows
	}
	return port.GetAllResponse{
		List: response,
	}, nil
}

func (m *Mock) GetByDate(ctx context.Context, date time.Time) (port.GetAllResponse, error) {
	response := []port.GetResponse{}
	for _, i := range m.resources {
		if i.Date == date {
			response = append(response, port.GetResponse{
				Id:        i.Id,
				Date:      i.Date,
				Amount:    i.Amount,
				PartnerId: i.PartnerId,
				TxRef:     i.TxRef,
				Status:    i.Status,
			})
		}
	}
	if len(response) == 0 {
		return port.GetAllResponse{}, port_commons.ErrSysNoRows
	}
	return port.GetAllResponse{
		List: response,
	}, nil
}

func (m *Mock) GetByTxRef(ctx context.Context, txRef string) (port.GetAllResponse, error) {
	response := []port.GetResponse{}
	for _, i := range m.resources {
		if i.TxRef == txRef {
			response = append(response, port.GetResponse{
				Id:        i.Id,
				Date:      i.Date,
				Amount:    i.Amount,
				PartnerId: i.PartnerId,
				TxRef:     i.TxRef,
				Status:    i.Status,
			})
		}
	}
	if len(response) == 0 {
		return port.GetAllResponse{}, port_commons.ErrSysNoRows
	}
	return port.GetAllResponse{
		List: response,
	}, nil
}

func (m *Mock) GetByStatus(ctx context.Context, status string) (port.GetAllResponse, error) {
	response := []port.GetResponse{}
	for _, i := range m.resources {
		if i.Status == status {
			response = append(response, port.GetResponse{
				Id:        i.Id,
				Date:      i.Date,
				Amount:    i.Amount,
				PartnerId: i.PartnerId,
				TxRef:     i.TxRef,
				Status:    i.Status,
			})
		}
	}
	if len(response) == 0 {
		return port.GetAllResponse{}, port_commons.ErrSysNoRows
	}
	return port.GetAllResponse{
		List: response,
	}, nil
}

func (m *Mock) UpdateByTransactionRef(ctx context.Context, req *port.UpdateByTransactionRefRequest) error {
	response := []MockTransaction{}
	for _, i := range m.resources {
		if i.TxRef == req.TransactionRef {
			response = append(response, MockTransaction{
				Id:        i.Id,
				Date:      i.Date,
				Amount:    i.Amount,
				PartnerId: i.PartnerId,
				TxRef:     i.TxRef,
				Status:    req.Status,
			})
		}
	}
	if len(response) == 0 {
		return port.ErrSysNoRows
	}
	m.resources = response
	return nil
}

func (m *Mock) UpdateStatus(ctx context.Context, req *port.UpdateRequest) error {
	response := []MockTransaction{}
	for _, i := range m.resources {
		if i.Id == req.Id {
			response = append(response, MockTransaction{
				Id:        i.Id,
				Date:      i.Date,
				Amount:    i.Amount,
				PartnerId: i.PartnerId,
				TxRef:     i.TxRef,
				Status:    req.Status,
			})
		}
	}
	if len(response) == 0 {
		return port_commons.ErrSysNoRows
	}
	m.resources = response
	return nil
}

func (m *Mock) GetByPartnerId(ctx context.Context, partnerId int) (port.GetAllResponse, error) {
	response := []port.GetResponse{}
	for _, i := range m.resources {
		if i.PartnerId == partnerId {
			response = append(response, port.GetResponse{
				Id:        i.Id,
				Date:      i.Date,
				Amount:    i.Amount,
				PartnerId: i.PartnerId,
				TxRef:     i.TxRef,
				Status:    i.Status,
			})
		}
	}
	if len(response) == 0 {
		return port.GetAllResponse{}, port_commons.ErrSysNoRows
	}
	return port.GetAllResponse{
		List: response,
	}, nil
}

func (m *Mock) Create(ctx context.Context, req *port.CreateRequest) (int, error) {
	newResourceId := len(m.resources) + 1
	m.resources = append(m.resources, MockTransaction{
		Id:        newResourceId,
		Amount:    req.Amount,
		PartnerId: req.PartnerId,
		Date:      time.Now(),
		TxRef:     req.TxRef,
		Status:    req.Status,
	})
	return newResourceId, nil
}

func (m *Mock) Delete(ctx context.Context, id int) error {
	updatedResources := []MockTransaction{}
	for _, i := range m.resources {
		if i.Id != id {
			updatedResources = append(updatedResources, MockTransaction{
				Id:        i.Id,
				Date:      i.Date,
				Amount:    i.Amount,
				PartnerId: i.PartnerId,
				TxRef:     i.TxRef,
				Status:    i.Status,
			})
		}
	}
	m.resources = updatedResources
	return nil
}

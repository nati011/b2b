package email

import (
	"context"
	"time"

	port "b2b.nati011.github.com/internal/port/application/payment/db"
	port_commons "b2b.nati011.github.com/internal/port/commons/db"
)

type Payment struct {
	Id             int
	OrderId        int
	PartnerId      int
	TransactionRef string
	Amount         float64
	Date           time.Time
}

type Mock struct {
	List []Payment
}

func NewMock() port.DB {
	return &Mock{}
}

func (p *Mock) Create(ctx context.Context, req port.CreateRequest) (int, error) {
	genId := len(p.List) + 1
	p.List = append(p.List, Payment{
		Id:             genId,
		OrderId:        req.OrderId,
		PartnerId:      req.PartnerId,
		TransactionRef: req.TransactionRef,
		Amount:         req.Amount,
		Date:           time.Now(),
	})
	return genId, nil
}

func (p *Mock) GetByID(ctx context.Context, id int) (port.GetResponse, error) {
	for _, m := range p.List {
		if m.Id == id {
			return port.GetResponse(m), nil
		}
	}
	return port.GetResponse{}, port_commons.ErrSysNoRows
}

func (p *Mock) GetByTransactionRef(ctx context.Context, transactionRef string) (port.GetResponse, error) {
	for _, m := range p.List {
		if m.TransactionRef == transactionRef {
			return port.GetResponse(m), nil
		}
	}
	return port.GetResponse{}, port_commons.ErrSysNoRows
}

func (p *Mock) GetByOrderId(ctx context.Context, orderId int) (port.GetAllResponse, error) {
	resp := port.GetAllResponse{}
	for _, m := range p.List {
		if m.OrderId == orderId {
			resp.List = append(resp.List, port.GetResponse(m))
		}
	}
	if len(resp.List) == 0 {
		return port.GetAllResponse{}, port_commons.ErrSysNoRows
	}
	return resp, nil
}

func (p *Mock) GetAll(ctx context.Context) (port.GetAllResponse, error) {
	resp := port.GetAllResponse{}
	for _, m := range p.List {
		resp.List = append(resp.List, port.GetResponse(m))
	}
	if len(resp.List) == 0 {
		return port.GetAllResponse{}, port_commons.ErrSysNoRows
	}
	return resp, nil
}

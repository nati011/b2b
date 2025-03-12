package invoice

import (
	"context"
	"math/rand"
	"time"

	port "b2b.nati011.github.com/internal/port/domain/invoice/db"
)

type Item struct {
	ProductId int
	Quantity  int
}

type MockInvoice struct {
	Id           int
	Created_Date time.Time
	ExternalId   string
	Status       string
	OrderId      int
	SubTotal     float64
	LineItems    []Item
	TaxAmount    float64
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
			items := []port.Item{}
			for _, i := range i.LineItems {
				items = append(items, port.Item{
					ProductId: i.ProductId,
					Quantity:  i.Quantity,
				})
			}
			return port.GetResponse{
				Id:           i.Id,
				Created_Date: i.Created_Date,
				ExternalId:   i.ExternalId,
				Status:       i.Status,
				OrderId:      i.OrderId,
				SubTotal:     i.SubTotal,
				LineItems:    items,
				TaxAmount:    i.TaxAmount,
			}, nil
		}
	}
	return port.GetResponse{}, port.ErrSysNoRows
}

func (m *Mock) GetAll(ctx context.Context) (port.GetAllResponse, error) {
	res := port.GetAllResponse{}
	for _, i := range m.invoices {
		items := []port.Item{}
		for _, i := range i.LineItems {
			items = append(items, port.Item{
				ProductId: i.ProductId,
				Quantity:  i.Quantity,
			})
		}
		res.List = append(res.List, port.GetResponse{
			Id:           i.Id,
			Created_Date: i.Created_Date,
			ExternalId:   i.ExternalId,
			Status:       i.Status,
			OrderId:      i.OrderId,
			SubTotal:     i.SubTotal,
			LineItems:    items,
			TaxAmount:    i.TaxAmount,
		})

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
			items := []port.Item{}
			for _, i := range i.LineItems {
				items = append(items, port.Item{
					ProductId: i.ProductId,
					Quantity:  i.Quantity,
				})
			}
			res.List = append(res.List, port.GetResponse{
				Id:           i.Id,
				Created_Date: i.Created_Date,
				ExternalId:   i.ExternalId,
				Status:       i.Status,
				OrderId:      i.OrderId,
				SubTotal:     i.SubTotal,
				LineItems:    items,
				TaxAmount:    i.TaxAmount,
			})
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
			items := []port.Item{}
			for _, i := range i.LineItems {
				items = append(items, port.Item{
					ProductId: i.ProductId,
					Quantity:  i.Quantity,
				})
			}
			res.List = append(res.List, port.GetResponse{
				Id:           i.Id,
				Created_Date: i.Created_Date,
				ExternalId:   i.ExternalId,
				Status:       i.Status,
				OrderId:      i.OrderId,
				SubTotal:     i.SubTotal,
				LineItems:    items,
				TaxAmount:    i.TaxAmount,
			})
		}
	}
	if len(res.List) == 0 {
		return port.GetAllResponse{}, port.ErrSysNoRows
	}
	return res, nil
}

func (m *Mock) GetByOrderId(ctx context.Context, orderId int) (port.GetResponse, error) {
	for _, i := range m.invoices {
		if i.OrderId == orderId {
			items := []port.Item{}
			for _, i := range i.LineItems {
				items = append(items, port.Item{
					ProductId: i.ProductId,
					Quantity:  i.Quantity,
				})
			}
			return port.GetResponse{
				Id:           i.Id,
				Created_Date: i.Created_Date,
				ExternalId:   i.ExternalId,
				Status:       i.Status,
				OrderId:      i.OrderId,
				SubTotal:     i.SubTotal,
				LineItems:    items,
				TaxAmount:    i.TaxAmount,
			}, nil
		}
	}
	return port.GetResponse{}, port.ErrSysNoRows
}

func (m *Mock) Create(ctx context.Context, req *port.CreateRequest) (int, error) {
	rand.Seed(time.Now().UnixNano())
	id := rand.Intn(1000-10+1) + 10
	mock_lineItems := []Item{}
	for _, i := range req.LineItems {
		mock_lineItems = append(mock_lineItems, Item{
			ProductId: i.ProductId,
			Quantity:  i.Quantity,
		})
	}
	m.invoices = append(m.invoices, MockInvoice{
		Id:           id,
		Created_Date: time.Now(),
		ExternalId:   req.ExternalId,
		Status:       req.Status,
		OrderId:      req.OrderId,
		SubTotal:     req.Subtotal,
		TaxAmount:    req.TaxAmount,
		LineItems:    mock_lineItems,
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
				OrderId:      i.OrderId,
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
				OrderId:      i.OrderId,
			})
		} else {
			res = append(res, MockInvoice(i))
		}
	}
	m.invoices = res
	return nil
}

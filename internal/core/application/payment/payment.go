package payment

import (
	"context"
	"errors"
	"time"

	port "b2b.nati011.github.com/internal/port/application/payment/db"
	port_commons "b2b.nati011.github.com/internal/port/commons/db"
)

var (
	ErrIdNotFound      = errors.New(" id not found")
	ErrEmptyGetContent = errors.New(" empty get content")
	ErrTxRefNotFound   = errors.New(" txRef not found")
	ErrOrderIdNotFound = errors.New(" orderId not found")
	ErrUnknown         = errors.New(" unknown error")
)

type GetAllResponse struct {
	List []GetResponse `json:"list"`
}

type GetResponse struct {
	Id             int       `json:"id"`
	OrderId        int       `json:"order_id"`
	PartnerId      int       `json:"partner_id"`
	TransactionRef string    `json:"tx_ref"`
	Amount         float64   `json:"amount"`
	Date           time.Time `json:"date"`
}

type CreateRequest struct {
	OrderId        int     `json:"order_id"`
	PartnerId      int     `json:"partner_id"`
	TransactionRef string  `json:"tx_ref"`
	Amount         float64 `json:"amount"`
}

type Provider interface {
	Create(context.Context, *CreateRequest) (int, error)
	GetByID(ctx context.Context, id int) (GetResponse, error)
	GetAll(ctx context.Context) (GetAllResponse, error)
	GetByTransactionRef(ctx context.Context, txRef string) (GetResponse, error)
	GetByOrderId(ctx context.Context, orderId int) (GetAllResponse, error)
}

type PaymentService struct {
	db port.DB
}

func NewPaymentService(DB port.DB) Provider {
	return &PaymentService{
		db: DB,
	}
}

func (p *PaymentService) Create(ctx context.Context, req *CreateRequest) (int, error) {
	id, err := p.db.Create(ctx, port.CreateRequest(*req))
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (p *PaymentService) GetByID(ctx context.Context, id int) (GetResponse, error) {
	resp, err := p.db.GetByID(ctx, id)
	if err != nil {
		return GetResponse{}, ErrIdNotFound
	}
	return GetResponse(resp), nil
}

func (p *PaymentService) GetAll(ctx context.Context) (GetAllResponse, error) {
	resp, err := p.db.GetAll(ctx)
	if err != nil {
		switch err {
		case port_commons.ErrSysNoRows:
			return GetAllResponse{}, ErrEmptyGetContent
		default:
			return GetAllResponse{}, ErrUnknown
		}
	}
	getAllResp := GetAllResponse{}
	for _, i := range resp.List {
		getAllResp.List = append(getAllResp.List, GetResponse(i))
	}
	return getAllResp, nil
}

func (p *PaymentService) GetByTransactionRef(ctx context.Context, txRef string) (GetResponse, error) {
	resp, err := p.db.GetByTransactionRef(ctx, txRef)
	if err != nil {
		return GetResponse{}, ErrTxRefNotFound
	}
	return GetResponse(resp), nil
}

func (p *PaymentService) GetByOrderId(ctx context.Context, orderId int) (GetAllResponse, error) {
	resp, err := p.db.GetByOrderId(ctx, orderId)
	if err != nil {
		switch err {
		case port_commons.ErrSysNoRows:
			return GetAllResponse{}, ErrOrderIdNotFound
		default:
			return GetAllResponse{}, ErrUnknown
		}
	}
	getAllResp := GetAllResponse{}
	for _, i := range resp.List {
		getAllResp.List = append(getAllResp.List, GetResponse(i))
	}
	return getAllResp, nil
}

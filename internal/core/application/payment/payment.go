package payment

import (
	"context"
	"time"

	port "b2b.nati011.github.com/internal/port/application/payment/db"
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
	Create(context.Context, CreateRequest) (int, error)
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

func (p *PaymentService) Create(context.Context, CreateRequest) (int, error) {
	return 0, nil
}

func (p *PaymentService) GetByID(ctx context.Context, id int) (GetResponse, error) {
	return GetResponse{}, nil
}

func (p *PaymentService) GetAll(ctx context.Context) (GetAllResponse, error) {
	return GetAllResponse{}, nil
}

func (p *PaymentService) GetByTransactionRef(ctx context.Context, txRef string) (GetResponse, error) {
	return GetResponse{}, nil
}

func (p *PaymentService) GetByOrderId(ctx context.Context, orderId int) (GetAllResponse, error) {
	return GetAllResponse{}, nil
}

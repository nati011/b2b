package payment

import (
	"context"
	"errors"
)

var (
	ErrUserNotSupplied           = errors.New("oopsy, user is mandatory")
	ErrAmountNotSupplied         = errors.New("oopsy, amount is mandatory")
	ErrAmountLessThanZero        = errors.New("oopsy, amount must be greater than zero")
	ErrPaymentPartnerNotSupplied = errors.New("oopsy, payment partner mandatory")
)

type CheckoutRequest struct {
	User_Id           int
	Amount            int64
	PaymentPartner_Id int
}

type CheckoutResponse struct {
	Checkout_url string
}

type Provider interface {
	Checkout(ctx context.Context, req *CheckoutRequest) (CheckoutResponse, error)
}

type PaymentService struct {
}

func NewPaymentService() Provider {
	return &PaymentService{}
}

func (p *PaymentService) Checkout(ctx context.Context, req *CheckoutRequest) (CheckoutResponse, error) {
	return CheckoutResponse{}, nil
}

package payment

import (
	"context"
	"errors"
)

var (
	ErrUserNotSupplied                 = errors.New("oopsy, user is mandatory")
	ErrAmountNotSupplied               = errors.New("oopsy, amount is mandatory")
	ErrAmountLessThanZero              = errors.New("oopsy, amount must be greater than zero")
	ErrPaymentPartnerNotSupplied       = errors.New("oopsy, payment partner id is mandatory")
	ErrTransactionReferenceNotSupplied = errors.New("oopsy, transaction refrence is mandatory")
)

type CheckoutRequest struct {
	User_Id           int
	Amount            int64
	PaymentPartner_Id int
}

type CheckoutResponse struct {
	Checkout_url string
}

type VerifyResponse struct {
	Is_verified bool
}

type Provider interface {
	Checkout(ctx context.Context, req *CheckoutRequest) (CheckoutResponse, error)
	Verify(ctx context.Context, tx_ref string) (VerifyResponse, error)
}

type PaymentService struct {
}

func NewPaymentService() Provider {
	return &PaymentService{}
}

func (p *PaymentService) Checkout(ctx context.Context, req *CheckoutRequest) (CheckoutResponse, error) {
	err := p.validateUserId(req.User_Id)
	if err != nil {
		return CheckoutResponse{}, err
	}
	err = p.validateAmount(req.Amount)
	if err != nil {
		return CheckoutResponse{}, err
	}
	err = p.validatePaymentPartner(req.PaymentPartner_Id)
	if err != nil {
		return CheckoutResponse{}, err
	}
	//gateway
	return CheckoutResponse{}, nil
}

func (p *PaymentService) Verify(ctx context.Context, tx_ref string) (VerifyResponse, error) {
	err := verifyTransactionRef(tx_ref)
	if err != nil {
		return VerifyResponse{}, err
	}
	return VerifyResponse{}, nil
}

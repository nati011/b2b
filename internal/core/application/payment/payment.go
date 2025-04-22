package payment

import (
	"context"
	"errors"
	"log"

	partner "b2b.nati011.github.com/internal/core/application/payment_partner"
	"b2b.nati011.github.com/internal/core/application/transaction"
	"b2b.nati011.github.com/internal/core/application/user"
	payment_processor "b2b.nati011.github.com/internal/core/domain/paymentProcessor"
)

var (
	ErrUserNotFound                    = errors.New("oopsy, user is found")
	ErrAmountNotSupplied               = errors.New("oopsy, amount is mandatory")
	ErrAmountLessThanZero              = errors.New("oopsy, amount must be greater than zero")
	ErrPaymentPartnerNotSupplied       = errors.New("oopsy, payment partner id is mandatory")
	ErrTransactionReferenceNotSupplied = errors.New("oopsy, transaction refrence is mandatory")
)

type CheckoutRequest struct {
	Amount           int64
	PaymentPartnerId int
}

type CheckoutResponse struct {
	Checkout_url string
}

type Provider interface {
	Checkout(ctx context.Context, req *CheckoutRequest) (CheckoutResponse, error)
	Verify(ctx context.Context, gateway_id int, tx_ref string) (bool, error)
	Callback(ctx context.Context, gateway_id int, tx_ref string)
}

type PaymentService struct {
	UserService        user.Provider
	PartnerService     partner.Provider
	TransactionService transaction.Provider
	PaymentProcessor   payment_processor.Provider
}

func NewPaymentService(
	ps partner.Provider,
	ts transaction.Provider) Provider {
	return &PaymentService{
		PartnerService:     ps,
		TransactionService: ts,
	}
}

func (p *PaymentService) Checkout(ctx context.Context, req *CheckoutRequest) (CheckoutResponse, error) {
	err := p.validateAmount(req.Amount)
	if err != nil {
		return CheckoutResponse{}, err
	}
	return CheckoutResponse{}, nil
}

func (p *PaymentService) Verify(ctx context.Context, gateway_id int, tx_ref string) (bool, error) {
	// part_resp, err := p.PartnerService.Get(ctx, gateway_id)

	return false, nil
}

func (p *PaymentService) Callback(ctx context.Context, gateway_id int, tx_ref string) {
	is_verified, err := p.Verify(ctx, gateway_id, tx_ref)
	if err != nil {
		switch err {
		default:
			log.Printf("failed to process incoming callback tx_ref: %v, gateway_id: %v", tx_ref, gateway_id)
		}
	}

	if is_verified {
		p.PaymentProcessor.Process(ctx, tx_ref)
	}
}

package payment

import (
	"context"
	"errors"
	"log"

	factory "b2b.nati011.github.com/internal/adapter/secondary/application/payment/gateway"
	partner "b2b.nati011.github.com/internal/core/application/payment_partner"
	"b2b.nati011.github.com/internal/core/application/transaction"
	payment "b2b.nati011.github.com/internal/port/application/payment/gateway"
)

var (
	ErrUserNotFound                    = errors.New("oopsy, user is found")
	ErrAmountNotSupplied               = errors.New("oopsy, amount is mandatory")
	ErrAmountLessThanZero              = errors.New("oopsy, amount must be greater than zero")
	ErrPaymentPartnerNotSupplied       = errors.New("oopsy, payment partner id is mandatory")
	ErrTransactionReferenceNotSupplied = errors.New("oopsy, transaction refrence is mandatory")
	ErrUnknown                         = errors.New("oopsy, unknown error has occured")
)

type CheckoutRequest struct {
	OrderId          int
	Amount           float64
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
	paymentPartner partner.Provider
	transaction    transaction.Provider
}

func NewPaymentService(partner partner.Provider, transaction transaction.Provider) Provider {
	return &PaymentService{
		paymentPartner: partner,
		transaction:    transaction,
	}
}

func (p *PaymentService) Checkout(ctx context.Context, req *CheckoutRequest) (CheckoutResponse, error) {
	err := p.validateAmount(req.Amount)
	if err != nil {
		return CheckoutResponse{}, err
	}
	paymentPartner, err := p.paymentPartner.GetPartnerSecret(ctx, req.PaymentPartnerId)

	if err != nil {
		return CheckoutResponse{}, err
	}

	paymentGateway, err := factory.PaymentPartnerFactory(paymentPartner.Name)

	if err != nil {
		return CheckoutResponse{}, err
	}

	paymentInitiateRequest := payment.InitiateRequest{
		Amount:         req.Amount,
		TransactionRef: req.OrderId,
		PartnerUrl:     paymentPartner.Init_payment_url,
		PartnerSecret:  paymentPartner.Secret,
	}

	checkoutUrl, err := paymentGateway.Initiate(paymentInitiateRequest)
	if err != nil {
		return CheckoutResponse{}, err
	}

	return CheckoutResponse{
		Checkout_url: checkoutUrl,
	}, nil
}

func (p *PaymentService) Verify(ctx context.Context, gateway_id int, tx_ref string) (bool, error) {
	return false, nil
}

func (p *PaymentService) Callback(ctx context.Context, gateway_id int, tx_ref string) {
	is_verified, err := p.Verify(ctx, gateway_id, tx_ref)
	if err != nil {
		switch err {
		default:
			log.Printf("failed to process incoming callback tx_ref: %v", tx_ref)
		}
	}

	if is_verified {
		//notify order
	}
}

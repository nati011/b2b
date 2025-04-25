package payment

import (
	"context"
	"errors"
	"log"

	factory "b2b.nati011.github.com/internal/adapter/secondary/application/payment/gateway"
	partner "b2b.nati011.github.com/internal/core/application/payment_partner"
	"b2b.nati011.github.com/internal/core/application/transaction"
	payment_processor "b2b.nati011.github.com/internal/core/domain/paymentProcessor"
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
	TransactionRef   string
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
	paymentPartner   partner.Provider
	transaction      transaction.Provider
	paymentProcessor payment_processor.Provider
	baseUrl          string
}

func NewPaymentService(partner partner.Provider, transaction transaction.Provider, processor payment_processor.Provider, baseUrl string) Provider {
	return &PaymentService{
		baseUrl:          baseUrl,
		paymentPartner:   partner,
		transaction:      transaction,
		paymentProcessor: processor,
	}
}

func (p *PaymentService) Checkout(ctx context.Context, req *CheckoutRequest) (CheckoutResponse, error) {
	order, err := p.paymentProcessor.FetchOrder(ctx, req.TransactionRef)
	paymentPartner, err := p.paymentPartner.GetPartnerSecret(ctx, req.PaymentPartnerId)
	if err != nil {
		return CheckoutResponse{}, ErrUnknown
	}

	paymentGateway, err := factory.PaymentPartnerFactory(paymentPartner.Name)
	if err != nil {
		return CheckoutResponse{}, err
	}

	paymentInitiateRequest := payment.InitiateRequest{
		Amount:         float64(order.Total),
		TransactionRef: req.TransactionRef,
		PartnerUrl:     paymentPartner.BaseUrl,
		PartnerSecret:  paymentPartner.Secret,
		BaseUrl:        p.baseUrl,
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
	paymentPartner, err := p.paymentPartner.GetPartnerSecret(ctx, gateway_id)
	if err != nil {
		return false, ErrUnknown
	}

	paymentGateway, err := factory.PaymentPartnerFactory(paymentPartner.Name)
	if err != nil {
		return false, err
	}

	paymentVerificationRequest := payment.VerificationRequest{
		PartnerUrl:     paymentPartner.BaseUrl,
		TransactionRef: tx_ref,
		PartnerSecret:  paymentPartner.Secret,
	}

	is_verified, err := paymentGateway.Verify(paymentVerificationRequest)
	if err != nil {
		return false, err
	}
	return is_verified, nil
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
		p.paymentProcessor.Process(ctx, tx_ref)
	}
}

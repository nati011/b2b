package payment

import (
	"context"
	"errors"
	"log"

	factory "b2b.nati011.github.com/internal/adapter/secondary/application/payment/gateway"
	partner "b2b.nati011.github.com/internal/core/application/payment_partner"
	"b2b.nati011.github.com/internal/core/application/transaction"
	payment_processor "b2b.nati011.github.com/internal/core/domain/paymentProcessor"
	db_port "b2b.nati011.github.com/internal/port/application/payment/db"
	gateway_port "b2b.nati011.github.com/internal/port/application/payment/gateway"
)

var (
	ErrUserNotFound                    = errors.New("oopsy, user is found")
	ErrAmountNotSupplied               = errors.New("oopsy, amount is mandatory")
	ErrAmountLessThanZero              = errors.New("oopsy, amount must be greater than zero")
	ErrPaymentPartnerNotSupported      = errors.New("oopsy, payment partner id is mandatory")
	ErrTransactionReferenceNotSupplied = errors.New("oopsy, transaction refrence is mandatory")
	ErrUnknown                         = errors.New("oopsy, unknown error has occured")
)

type CheckoutRequest struct {
	Amount           float64
	PaymentPartnerId int
	TransactionRef   string
}

type CheckoutResponse struct {
	Checkout_url string
}

type Provider interface {
	Checkout(ctx context.Context, req *CheckoutRequest) (CheckoutResponse, error)
	Verify(ctx context.Context, PaymentPartnerId int, tx_ref string) (bool, error)
	Callback(ctx context.Context, PaymentPartnerId int, tx_ref string)
}

type PaymentService struct {
	DB                 db_port.DB
	PartnerService     partner.Provider
	TransactionService transaction.Provider
	PaymentProcessor   payment_processor.Provider
}

func NewPaymentService(
	partner partner.Provider,
	transaction transaction.Provider,
	paymentProcessor payment_processor.Provider) Provider {
	return &PaymentService{
		PartnerService:     partner,
		TransactionService: transaction,
		PaymentProcessor:   paymentProcessor,
	}
}

func (p *PaymentService) Checkout(ctx context.Context, req *CheckoutRequest) (CheckoutResponse, error) {
	if req.TransactionRef == "" {
		return CheckoutResponse{}, ErrTransactionReferenceNotSupplied
	}
	paymentPartner, err := p.PartnerService.GetPartnerSecret(ctx, req.PaymentPartnerId)
	if err != nil {
		switch err {
		case partner.ErrIdNotFound:
			return CheckoutResponse{}, ErrPaymentPartnerNotSupported
		default:
			return CheckoutResponse{}, ErrUnknown
		}
	}

	paymentGateway, err := factory.PaymentPartnerFactory(paymentPartner.Name)
	if err != nil {
		return CheckoutResponse{}, err
	}

	checkoutUrl, err := paymentGateway.Initiate(gateway_port.InitiateRequest{
		Amount:         req.Amount,
		TransactionRef: req.TransactionRef,
		PartnerUrl:     paymentPartner.BaseURL,
		PartnerSecret:  paymentPartner.Secret,
	})
	if err != nil {
		return CheckoutResponse{}, err
	}

	// create payment process

	return CheckoutResponse{
		Checkout_url: checkoutUrl,
	}, nil
}

func (p *PaymentService) Verify(ctx context.Context, gateway_id int, tx_ref string) (bool, error) {
	if tx_ref == "" {
		return false, ErrTransactionReferenceNotSupplied
	}

	paymentPartner, err := p.PartnerService.GetPartnerSecret(ctx, gateway_id)
	if err != nil {
		switch err {
		case partner.ErrIdNotFound:
			return false, ErrPaymentPartnerNotSupported
		default:
			return false, ErrUnknown
		}
	}

	paymentGateway, err := factory.PaymentPartnerFactory(paymentPartner.Name)
	if err != nil {
		return false, err
	}

	paymentVerificationRequest := gateway_port.VerificationRequest{
		PartnerUrl:     paymentPartner.BaseURL,
		TransactionRef: tx_ref,
		PartnerSecret:  paymentPartner.Secret,
	}

	is_verified, err := paymentGateway.Verify(paymentVerificationRequest)
	if err != nil {
		return false, ErrUnknown
	}

	// // if verified create transaction
	// p.TransactionService.Create(ctx, &transaction.CreateRequest{
	// 	Amount: ,
	// 	Partner_Id:
	// })
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
		p.PaymentProcessor.Process(ctx, tx_ref)
	}
}

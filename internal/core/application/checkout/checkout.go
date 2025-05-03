package checkout

import (
	"context"
	"errors"
	"time"

	factory "b2b.nati011.github.com/internal/adapter/secondary/application/payment/gateway"
	partner "b2b.nati011.github.com/internal/core/application/payment_partner"
	"b2b.nati011.github.com/internal/core/application/transaction"
	port "b2b.nati011.github.com/internal/port/application/payment/db"
	payment "b2b.nati011.github.com/internal/port/application/payment/gateway"
)

var (
	ErrUserNotFound                    = errors.New("oopsy, user is found")
	ErrAmountNotSupplied               = errors.New("oopsy, amount is mandatory")
	ErrAmountLessThanZero              = errors.New("oopsy, amount must be greater than zero")
	ErrPaymentPartnerNotSupported      = errors.New("oopsy, payment partner id is not supported")
	ErrTransactionReferenceNotSupplied = errors.New("oopsy, transaction refrence is mandatory")
	ErrOrderIdNotSupplied              = errors.New("oopsy, order id not supplied")
	ErrUnknown                         = errors.New("oopsy, unknown error has occured")
)

type CheckoutRequest struct {
	OrderId          int
	Amount           float64
	PaymentPartnerId int
}

type CheckoutResponse struct {
	CheckoutUrl    string
	TransactionRef string
}

type Provider interface {
	Checkout(ctx context.Context, req *CheckoutRequest) (CheckoutResponse, error)
}

type CheckoutService struct {
	db             port.DB
	paymentPartner partner.Provider
	transaction    transaction.Provider
	frontendUrl    string
	baseUrl        string
}

type CreatePaymentRequest struct {
	PaymentPartnerId int
	OrderId          int
	Amount           float64
}

func NewCheckoutService(DB port.DB, partner partner.Provider, transaction transaction.Provider, frontendUrl string, baseUrl string) Provider {
	return &CheckoutService{
		db:             DB,
		frontendUrl:    frontendUrl,
		baseUrl:        baseUrl,
		paymentPartner: partner,
		transaction:    transaction,
	}
}

func (p *CheckoutService) CreatePayment(ctx context.Context, req *CreatePaymentRequest) (string, error) {
	currentTimestamp := time.Now()
	generatedTxRef := currentTimestamp.Format("2006_01_02_15_04_05")
	createRequest := port.CreateRequest{
		OrderId:        req.OrderId,
		PartnerId:      req.PaymentPartnerId,
		TransactionRef: generatedTxRef,
	}
	_, err := p.transaction.Create(ctx, &transaction.CreateRequest{
		Amount:    req.Amount,
		PartnerId: req.PaymentPartnerId,
		TxRef:     generatedTxRef,
		Status:    transaction.PENDING_STATUS,
	})

	_, err = p.db.Create(ctx, createRequest)
	if err != nil {
		return "", err
	}
	return generatedTxRef, err
}

func (p *CheckoutService) Checkout(ctx context.Context, req *CheckoutRequest) (CheckoutResponse, error) {
	paymentPartner, err := p.paymentPartner.GetPartnerSecret(ctx, req.PaymentPartnerId)
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

	transaction_ref, err := p.CreatePayment(ctx, &CreatePaymentRequest{
		PaymentPartnerId: req.PaymentPartnerId,
		OrderId:          req.OrderId,
	})
	if err != nil {
		return CheckoutResponse{}, err
	}

	_, err = p.transaction.Create(ctx, &transaction.CreateRequest{
		Amount:    req.Amount,
		PartnerId: req.PaymentPartnerId,
		TxRef:     transaction_ref,
		Status:    transaction.PENDING_STATUS,
	})
	if err != nil {
		return CheckoutResponse{}, err
	}

	paymentInitiateRequest := payment.InitiateRequest{
		Amount:         req.Amount,
		TransactionRef: transaction_ref,
		PartnerUrl:     paymentPartner.BaseURL,
		PartnerSecret:  paymentPartner.Secret,
		BaseUrl:        p.baseUrl,
		ReturnUrl:      p.frontendUrl,
	}

	checkoutUrl, err := paymentGateway.Initiate(paymentInitiateRequest)
	if err != nil {
		return CheckoutResponse{}, err
	}

	return CheckoutResponse{
		TransactionRef: transaction_ref,
		CheckoutUrl:    checkoutUrl,
	}, nil
}

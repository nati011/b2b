package checkout

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	factory "b2b.nati011.github.com/internal/adapter/secondary/application/payment/gateway"
	partner "b2b.nati011.github.com/internal/core/application/payment_partner"
	"b2b.nati011.github.com/internal/core/application/transaction"
	port "b2b.nati011.github.com/internal/port/application/payment/db"
	payment "b2b.nati011.github.com/internal/port/application/payment/gateway"
	"github.com/google/uuid"
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
	//to ensure uniqueness use current time stamp as transaction ref
	generatedTxRef := fmt.Sprintf("%s_%s", currentTimestamp.Format("2006_01_02_15_04_05"), uuid.New().String())
	createRequest := port.CreateRequest{
		OrderId:        req.OrderId,
		PartnerId:      req.PaymentPartnerId,
		TransactionRef: generatedTxRef,
	}

	_, err := p.db.Create(ctx, createRequest)
	if err != nil {
		return "", err
	}
	return generatedTxRef, err
}

func (p *CheckoutService) Checkout(ctx context.Context, req *CheckoutRequest) (CheckoutResponse, error) {
	if req.Amount == 0 {
		return CheckoutResponse{}, ErrAmountNotSupplied
	}
	paymentPartner, err := p.paymentPartner.Get(ctx, req.PaymentPartnerId)
	if err != nil {
		switch err {
		case partner.ErrIdNotFound:
			return CheckoutResponse{}, ErrPaymentPartnerNotSupported
		default:
			return CheckoutResponse{}, ErrUnknown
		}
	}
	if paymentPartner.Status != partner.ACTIVE_STATUS {
		return CheckoutResponse{}, ErrPaymentPartnerNotSupported
	}

	paymentPartnerSecret, err := p.paymentPartner.GetPartnerSecret(ctx, req.PaymentPartnerId)
	if err != nil {
		log.Printf("Error while fetching paymentPartnerSecret: %v", err.Error())
		switch err {
		case partner.ErrIdNotFound:
			return CheckoutResponse{}, ErrPaymentPartnerNotSupported
		default:
			return CheckoutResponse{}, ErrUnknown
		}
	}

	paymentGateway, err := factory.PaymentPartnerFactory(paymentPartner.Name)
	if err != nil {
		return CheckoutResponse{}, ErrPaymentPartnerNotSupported
	}

	transaction_ref, err := p.CreatePayment(ctx, &CreatePaymentRequest{
		Amount:           req.Amount,
		PaymentPartnerId: req.PaymentPartnerId,
		OrderId:          req.OrderId,
	})
	if err != nil {
		log.Printf("Error while fetching creating payment: %v", err.Error())
		return CheckoutResponse{}, err
	}

	returnUrl := fmt.Sprintf("%v/%v", p.frontendUrl, req.OrderId)

	paymentInitiateRequest := payment.InitiateRequest{
		Amount:         req.Amount,
		TransactionRef: transaction_ref,
		PartnerUrl:     paymentPartner.BaseURL,
		PartnerSecret:  paymentPartnerSecret.Secret,
		BaseUrl:        p.baseUrl,
		ReturnUrl:      returnUrl,
	}

	checkoutUrl, err := paymentGateway.Initiate(paymentInitiateRequest)
	if err != nil {
		log.Printf("Error while fetching Initiating payment: %v", err.Error())
		return CheckoutResponse{}, err
	}

	_, err = p.transaction.Create(ctx, &transaction.CreateRequest{
		Amount:    req.Amount,
		PartnerId: req.PaymentPartnerId,
		TxRef:     transaction_ref,
		Status:    transaction.PENDING_STATUS,
	})
	if err != nil {
		log.Printf("Error while creating transaction: %v", err.Error())
		switch err {
		case transaction.ErrAmountIsNotSupplied:
			return CheckoutResponse{}, ErrAmountNotSupplied
		default:
			return CheckoutResponse{}, err
		}
	}

	return CheckoutResponse{
		TransactionRef: transaction_ref,
		CheckoutUrl:    checkoutUrl,
	}, nil
}

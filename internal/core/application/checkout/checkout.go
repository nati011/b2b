package checkout

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	factory "b2b.nati011.github.com/internal/adapter/secondary/application/payment/gateway"
	payment "b2b.nati011.github.com/internal/core/application/payment"
	partner "b2b.nati011.github.com/internal/core/application/payment_partner"
	"b2b.nati011.github.com/internal/core/application/transaction"
	payment_gateway "b2b.nati011.github.com/internal/port/application/payment/gateway"
	"github.com/google/uuid"
)

var (
	ErrUserNotFound                    = errors.New(" user is found")
	ErrAmountNotSupplied               = errors.New(" amount is mandatory")
	ErrAmountLessThanZero              = errors.New(" amount must be greater than zero")
	ErrPaymentPartnerNotSupported      = errors.New(" payment partner id is not supported")
	ErrTransactionReferenceNotSupplied = errors.New(" transaction refrence is mandatory")
	ErrOrderIdNotSupplied              = errors.New(" order id not supplied")
	ErrInvalidPaymentMethod            = errors.New(" invalid payment method")
	ErrUnknown                         = errors.New(" unknown error has occured")
)

const (
	PAYMENT_METHOD_DIGITAL = "DIGITAL_PAYMENT"
	PAYMENT_METHOD_MANUAL  = "MANUAL_PAYMENT"
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

type ReinitiateCheckoutRequest struct {
	OrderId int
}

type SubscriptionPaymentRequest struct {
	SubscriptionId   int
	Amount           float64
	PaymentPartnerId int
}

type SubscriptionPaymentResponse struct {
	CheckoutUrl    string
	TransactionRef string
}

type ReinitiateSubscriptionPaymentRequest struct {
	SubscriptionId int
}

// TODO: generalize naming to address expanded domain responsibility as central payment hub
type Provider interface {
	Checkout(ctx context.Context, req *CheckoutRequest) (CheckoutResponse, error)
	ReinitiateCheckout(ctx context.Context, req *ReinitiateCheckoutRequest) (CheckoutResponse, error)
	SubscriptionPayment(ctx context.Context, req *SubscriptionPaymentRequest) (SubscriptionPaymentResponse, error)
	ReinitiateSubscriptionPayment(ctx context.Context, req *ReinitiateSubscriptionPaymentRequest) (SubscriptionPaymentResponse, error)
}

type CheckoutService struct {
	paymentService payment.Provider
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

func NewCheckoutService(paymentService payment.Provider, partner partner.Provider, transaction transaction.Provider, frontendUrl string, baseUrl string) Provider {
	return &CheckoutService{
		paymentService: paymentService,
		frontendUrl:    frontendUrl,
		baseUrl:        baseUrl,
		paymentPartner: partner,
		transaction:    transaction,
	}
}

func (p *CheckoutService) CreatePayment(ctx context.Context, req *CreatePaymentRequest) (string, error) {
	currentTimestamp := time.Now()
	//to ensure uniqueness use current time stamp as transaction ref
	generatedTxRef := fmt.Sprintf("%s_%s", currentTimestamp.Format("Monday_15"), uuid.New().String())
	_, err := p.paymentService.Create(ctx, &payment.CreateRequest{
		OrderId:        req.OrderId,
		PartnerId:      req.PaymentPartnerId,
		TransactionRef: generatedTxRef,
		Amount:         req.Amount,
	})
	if err != nil {
		return "", err
	}
	return generatedTxRef, err
}

func (p *CheckoutService) Checkout(ctx context.Context, req *CheckoutRequest) (CheckoutResponse, error) {
	log.Printf("checkout for orderId: %v", req.OrderId)
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
	switch paymentPartner.PaymentMethod {
	case PAYMENT_METHOD_MANUAL:
		transaction_ref, err := p.CreatePayment(ctx, &CreatePaymentRequest{
			Amount:           req.Amount,
			PaymentPartnerId: req.PaymentPartnerId,
			OrderId:          req.OrderId,
		})
		if err != nil {
			log.Printf("Error while creating payment: %v", err.Error())
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
		}, nil
	case PAYMENT_METHOD_DIGITAL:
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
			log.Printf("Error while creating payment: %v", err.Error())
			return CheckoutResponse{}, err
		}

		returnUrl := fmt.Sprintf("%v/%v", p.frontendUrl, req.OrderId)

		paymentInitiateRequest := payment_gateway.InitiateRequest{
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

	return CheckoutResponse{}, ErrInvalidPaymentMethod
}

func (p *CheckoutService) ReinitiateCheckout(ctx context.Context, req *ReinitiateCheckoutRequest) (CheckoutResponse, error) {
	log.Printf("reinitate payment for orderId: %v", req.OrderId)
	payRecords, err := p.paymentService.GetByOrderId(ctx, req.OrderId)
	if err != nil {
		log.Printf("failed to fetch payment records: %v", err)
		return CheckoutResponse{}, ErrUnknown
	}

	pay := payRecords.List[len(payRecords.List)-1]
	paymentPartner, err := p.paymentPartner.Get(ctx, pay.PartnerId)
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

	paymentPartnerSecret, err := p.paymentPartner.GetPartnerSecret(ctx, paymentPartner.Id)
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
		Amount:           pay.Amount,
		PaymentPartnerId: pay.PartnerId,
		OrderId:          pay.OrderId,
	})
	if err != nil {
		log.Printf("Error while creating payment: %v", err.Error())
		return CheckoutResponse{}, err
	}

	paymentInitiateRequest := payment_gateway.InitiateRequest{
		Amount:         pay.Amount,
		TransactionRef: transaction_ref,
		PartnerUrl:     paymentPartner.BaseURL,
		PartnerSecret:  paymentPartnerSecret.Secret,
		BaseUrl:        p.baseUrl,
		ReturnUrl:      p.frontendUrl,
	}

	checkoutUrl, err := paymentGateway.Initiate(paymentInitiateRequest)
	if err != nil {
		log.Printf("Error while fetching Initiating payment: %v", err.Error())
		return CheckoutResponse{}, err
	}

	_, err = p.transaction.Create(ctx, &transaction.CreateRequest{
		Amount:    pay.Amount,
		PartnerId: pay.PartnerId,
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

func (p *CheckoutService) SubscriptionPayment(ctx context.Context, req *SubscriptionPaymentRequest) (SubscriptionPaymentResponse, error) {
	return SubscriptionPaymentResponse{}, nil
}

func (p *CheckoutService) ReinitiateSubscriptionPayment(ctx context.Context, req *ReinitiateSubscriptionPaymentRequest) (SubscriptionPaymentResponse, error) {
	return SubscriptionPaymentResponse{}, nil
}
